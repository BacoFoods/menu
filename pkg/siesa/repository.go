package siesa

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/siesa/repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(reference *Reference) error {
	if err := r.db.Save(reference).Error; err != nil {
		shared.LogError("error creating reference", LogDBRepository, "Create", err, reference)
		return err
	}
	return nil
}

// TruncateRecords truncates all records
// TruncateRecords truncates all records in the YourModel table
func (r *DBRepository) TruncateRecords() error {
	// Comienza una transacción
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Trunca la tabla 'references' dentro de la transacción con condiciones
	if err := tx.Unscoped().Where("true").Delete(&Reference{}).Error; err != nil {
		// Revierte la transacción en caso de error
		tx.Rollback()
		return err
	}

	// Confirma la transacción si la eliminación fue exitosa
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *DBRepository) Find(filters map[string]string) (*Reference, error) {
	var reference Reference
	if err := r.db.Unscoped().Preload(clause.Associations).Find(&reference, filters).Limit(1).Error; err != nil {
		shared.LogError("error getting reference row", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return &reference, nil
}

func (r *DBRepository) FindReferences(filters map[string]string) ([]Reference, error) {
	var reference []Reference
	if err := r.db.Find(&reference, filters).Error; err != nil {
		shared.LogError("error finding reference", LogDBRepository, "FindReferences", err, filters)
		return nil, err
	}
	return reference, nil
}

// Create method for create a new category in database
func (r *DBRepository) CreateReference(reference *Reference) (*Reference, error) {
	if err := r.db.Preload(clause.Associations).Save(reference).Error; err != nil {
		shared.LogError("error creating category", LogDBRepository, "CreateReference", err, reference)
		return nil, err
	}
	return reference, nil
}

// Delete method for delete a reference in database
func (r *DBRepository) DeleteReference(referenceID string) (*Reference, error) {
	var reference Reference
	if err := r.db.First(&reference, referenceID).Error; err != nil {
		shared.LogError("error getting category", LogDBRepository, "Delete", err, referenceID)
		return nil, err
	}
	if err := r.db.Delete(&reference).Error; err != nil {
		shared.LogError("error deleting category", LogDBRepository, "Delete", err, reference)
		return nil, err
	}
	return &reference, nil
}

// Update method for update a reference in database
func (r *DBRepository) UpdateReference(reference *Reference) (*Reference, error) {
	var referenceDB Reference
	if err := r.db.Preload(clause.Associations).First(&referenceDB, reference.ID).Error; err != nil {
		shared.LogError("error getting category", LogDBRepository, "UpdateReference", err, reference)
		return nil, err
	}
	if err := r.db.Model(&referenceDB).Updates(reference).Error; err != nil {
		shared.LogError("error updating category", LogDBRepository, "UpdateReference", err, reference)
		return nil, err
	}
	return &referenceDB, nil
}
