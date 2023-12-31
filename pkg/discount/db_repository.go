package discount

import (
	"fmt"
	"strings"

	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/discount/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(discount *Discount) (*Discount, error) {
	if err := r.db.Save(discount).Error; err != nil {
		shared.LogError("error creating discount", LogDBRepository, "Create", err, discount)
		return nil, err
	}
	return discount, nil
}

func (r *DBRepository) GetMany(ids []uint) ([]Discount, error) {
	var discount []Discount

	if len(ids) == 0 {
		return discount, nil
	}

	if err := r.db.Where("id IN ?", ids).Find(&discount).Error; err != nil {
		shared.LogError("error getting discounts", LogDBRepository, "GetMany", err, ids)
		return nil, err
	}

	return discount, nil
}

func (r *DBRepository) Find(filter map[string]string) ([]Discount, error) {
	var discount []Discount
	if err := r.db.Find(&discount, filter).Error; err != nil {
		shared.LogError("error finding discount", LogDBRepository, "Find", err)
		return nil, err
	}
	return discount, nil
}

func (r *DBRepository) Get(discountID string) (*Discount, error) {
	if strings.TrimSpace(discountID) == "" {
		err := fmt.Errorf(ErrorDiscountIDEmpty)
		shared.LogWarn("error getting discount", LogDBRepository, "Get", err)
		return nil, err
	}

	var discount Discount
	if err := r.db.Preload(clause.Associations).First(&discount, discountID).Error; err != nil {
		shared.LogError("error getting discount", LogDBRepository, "Find", err, discountID)
		return nil, err
	}
	return &discount, nil
}

func (r *DBRepository) Update(discount Discount) (*Discount, error) {
	var discountDB Discount
	if err := r.db.First(&discountDB, discount.ID).Error; err != nil {
		shared.LogError("error getting discount", LogDBRepository, "Update", err, discount.ID, discount)
		return nil, err
	}

	if err := r.db.Model(&discountDB).Updates(discount).Error; err != nil {
		shared.LogError("error updating discount", LogDBRepository, "Update", err, discount.ID, discount, discount)
		return nil, err
	}
	return &discountDB, nil
}

func (r *DBRepository) Delete(discountID string) (*Discount, error) {
	var discountDB Discount
	if err := r.db.First(&discountDB, discountID).Error; err != nil {
		shared.LogError("error getting discount", LogDBRepository, "Delete", err, discountID)
		return nil, err
	}

	if err := r.db.Delete(&discountDB, discountID).Error; err != nil {
		shared.LogError("error deleting discount", LogDBRepository, "Delete", err, discountDB)
		return nil, err
	}
	return &discountDB, nil
}
