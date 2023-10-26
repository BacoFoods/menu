package tables

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	LogService = "pkg/tables/service"
)

type Service interface {
	Get(id string) (*Table, error)
	Find(query map[string]any) ([]Table, error)
	Create(table *Table) (*Table, error)
	Update(id string, table *Table) (*Table, error)
	Delete(id string) error
	ScanQR(qrID string) (*Table, error)
	GenerateQR(tableId string) (*Table, error)
}

type service struct {
	repository Repository
	oitHost    string
}

func NewService(repository Repository, oitHost string) service {
	return service{repository, oitHost}
}

func (s service) Get(id string) (*Table, error) {
	return s.repository.Get(id)
}

func (s service) Find(query map[string]any) ([]Table, error) {
	return s.repository.Find(query)
}

func (s service) Create(table *Table) (*Table, error) {
	return s.repository.Create(table)
}

func (s service) Update(id string, table *Table) (*Table, error) {
	return s.repository.Update(id, table)
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s service) ScanQR(qrID string) (*Table, error) {
	return s.repository.ScanQR(qrID)
}

func (s service) GenerateQR(tableId string) (*Table, error) {
	table, err := s.Get(tableId)
	if err != nil {
		return nil, err
	}

	if table.QR != nil {
		return nil, errors.New("table already has a qr")
	}

	now := time.Now()
	qrID := uuid.New().String()
	qr := QR{
		TableID:   table.ID,
		IsActive:  true,
		DisplayID: qrID,
		CreatedAt: &now,
		UpdatedAt: &now,
		URL:       s.buildQRURL(qrID),
	}

	nQR, err := s.repository.CreateQR(qr)
	if err != nil {
		return nil, err
	}

	table.QR = nQR

	return table, nil
}

func (s service) buildQRURL(qrID string) string {
	return fmt.Sprintf("%s/%s", s.oitHost, qrID)
}
