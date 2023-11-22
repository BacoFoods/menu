package client

import (
	"time"

	"gorm.io/gorm"
)

const (
	ErrorClientCreating = "error creating client"
	ErrorClientUpdating = "error updating client"
	ErrorClientDeleting = "error deleting client"
	ErrorClientListing  = "error listing clients"
	ErrorClientGetting  = "error getting client"
	ErrorClientIDEmpty  = "error client id empty"
)

type Repository interface {
	Create(client *Client) (*Client, error)
	Update(client *Client) (*Client, error)
	Delete(id string) (*Client, error)
	List(map[string]any) ([]Client, error)
	Get(id string) (*Client, error)
}

type Client struct {
	ID           uint            `json:"id"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	DocumentType string          `json:"document_type"`
	Document     string          `json:"document"`
	Address      string          `json:"address"`
	CustomerID   *string         `json:"customer_id"`
	CreatedAt    *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt    *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt    *gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
}

func DefaultClient() *Client {
	return &Client{
		Name:     "consumidor final",
		Document: "222222222222",
	}
}
