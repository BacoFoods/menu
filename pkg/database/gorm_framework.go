package database

import (
	"errors"
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/sirupsen/logrus"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Many2ManyEntity interface {
	JoinTable(db gorm.DB) error
}

type GormFramework struct {
	dbConfig internal.DBConfig
	db       *gorm.DB
}

func MustNewGormFramework(customDSN string) *GormFramework {
	if strings.TrimSpace(customDSN) == "" {
		customDSN = dsn()
	}

	logrus.Debugf("Connecting to database: %s", escapeDSN(customDSN))
	db, err := gorm.Open(gormpostgres.Open(dsn()), &gorm.Config{
		TranslateError: true,
		NowFunc: func() time.Time {
			return time.Now().UTC().Add(-5 * time.Hour)
		},
	})
	if err != nil {
		logrus.Fatalf("error connecting to database: %s", escapeDSN(customDSN))
		return nil
	}

	return &GormFramework{
		dbConfig: internal.Config.DBConfig,
		db:       db,
	}
}

func (g *GormFramework) GetDBClient() *gorm.DB {
	return g.db
}

func (g *GormFramework) MustMakeMigrations(entities ...any) {
	for _, entity := range entities {
		if err := g.db.AutoMigrate(entity); err != nil {
			logrus.Fatal(fmt.Sprintf("error migrating %+v error:%s", entity, err.Error()))
		}

		switch entityType := entity.(type) {
		case Many2ManyEntity:
			if err := entityType.JoinTable(*g.db); err != nil {
				logrus.Fatal(fmt.Sprintf("error joining table %+v error:%s", entity, err.Error()))
			}
		}
	}
}

func (g *GormFramework) MustSetSeeds(entity any, seeds ...any) {
	for _, seed := range seeds {
		if g.db.Migrator().HasTable(entity) {
			if err := g.db.First(entity).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				if err := g.db.Create(seed).Error; err != nil {
					logrus.Fatal(fmt.Sprintf("error migrating seed %+v error:%s", seed, err.Error()))
				}
			}
		}
	}
}

func dsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		internal.Config.DBConfig.Host,
		internal.Config.DBConfig.User,
		internal.Config.DBConfig.Password,
		internal.Config.DBConfig.Name,
		internal.Config.DBConfig.Port,
	)
}

func escapeDSN(customDSN string) string {
	if strings.TrimSpace(customDSN) == "" {
		customDSN = dsn()
	}
	return strings.ReplaceAll(
		customDSN,
		internal.Config.DBConfig.Password,
		"***",
	)
}
