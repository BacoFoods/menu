package connector

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/connector/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(equivalence *Equivalence) (*Equivalence, error) {
	if err := r.db.Save(equivalence).Error; err != nil {
		shared.LogError("error creating country", LogDBRepository, "Create", err, equivalence)
		return nil, err
	}

	return equivalence, nil
}

func (r *DBRepository) FindReference(filter map[string]string) (*Equivalence, error) {
	tx := r.db.
		Preload(clause.Associations)

	// Handle specific filters for ChannelID and ProductID
	if channelID, ok := filter["channel_id"]; ok {
		tx = tx.Where("channel_id = ?", channelID)
		delete(filter, "channel_id")
	}

	if productID, ok := filter["product_id"]; ok {
		tx = tx.Where("product_id = ?", productID)
		delete(filter, "product_id")
	}

	var equivalence Equivalence
	if err := tx.Limit(1).Find(&equivalence, filter).Error; err != nil {
		shared.LogError("error finding equivalences", LogDBRepository, "FindReference", err, filter)
		return nil, err
	}

	return &equivalence, nil
}
