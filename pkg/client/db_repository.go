package client

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const (
	LogRepository string = "pkg/client/db_repository"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(client *Client) (*Client, error) {
	if err := r.db.Create(client).Error; err != nil {
		shared.LogError("error creating client", LogRepository, "Create", err, *client)
		return nil, err
	}

	return client, nil
}

func (r *DBRepository) Update(client *Client) (*Client, error) {
	if err := r.db.Save(client).Error; err != nil {
		shared.LogError("error updating client", LogRepository, "Update", err, *client)
		return nil, err
	}

	return client, nil
}

func (r *DBRepository) Delete(id string) (*Client, error) {
	var client Client
	if err := r.db.First(&client, id).Error; err != nil {
		shared.LogError("error getting client", LogRepository, "Delete", err, id)
		return nil, err
	}

	if err := r.db.Delete(&client).Error; err != nil {
		shared.LogError("error deleting client", LogRepository, "Delete", err, client)
		return nil, err
	}

	return &client, nil
}

func (r *DBRepository) List() ([]Client, error) {
	var clients []Client
	if err := r.db.Find(&clients).Error; err != nil {
		shared.LogError("error listing clients", LogRepository, "List", err)
		return nil, err
	}

	return clients, nil
}

func (r *DBRepository) Get(id string) (*Client, error) {
	if id == "" {
		shared.LogWarn("error getting client", LogRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var client Client
	if err := r.db.First(&client, id).Error; err != nil {
		shared.LogError("error getting client", LogRepository, "Get", err, id)
		return nil, err
	}

	return &client, nil
}
