package tables

import (
	"errors"
	"fmt"
	"time"

	"github.com/BacoFoods/menu/pkg/shared"
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

	FindZones(filters map[string]any) ([]Zone, error)
	GetZone(zoneID string) (*Zone, error)
	CreateZone(zone *Zone, tableNumber, tableAmount int) (*Zone, error)
	UpdateZone(zonID string, zone *Zone) (*Zone, error)
	DeleteZone(zoneID string) error
	AddTables(zoneID string, tables []uint) error
	RemoveTables(zoneID string, tables []uint) error
	EnableZone(zoneID string) (*Zone, error)
}

type service struct {
	zones      *zoneRepository
	repository *tableRepository
	oitHost    string
}

func NewService(repository *tableRepository, zones *zoneRepository, oitHost string) service {
	return service{zones, repository, oitHost}
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
	table, err := s.repository.ScanQR(qrID)
	if err != nil {
		return nil, err
	}

	if table == nil || table.Zone == nil || table.Zone.StoreID == nil {
		return nil, nil
	}

	return table, nil
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
		TableID:   &table.ID,
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

func (s service) FindZones(filters map[string]any) ([]Zone, error) {
	return s.zones.Find(filters)
}

func (s service) GetZone(zoneID string) (*Zone, error) {
	return s.zones.GetZone(zoneID)
}

func (s service) CreateZone(zone *Zone, tableNumber, tableAmount int) (*Zone, error) {
	SetTables(zone, tableNumber, tableAmount)
	return s.zones.Create(zone)
}

func (s service) UpdateZone(zoneID string, zone *Zone) (*Zone, error) {
	return s.zones.Update(zoneID, zone)
}

func (s service) DeleteZone(zoneID string) error {
	return s.zones.Delete(zoneID)
}

func (s service) AddTables(zoneID string, tables []uint) error {
	zone, err := s.zones.GetZone(zoneID)
	if err != nil {
		return err
	}

	if zone == nil {
		err := fmt.Errorf(ErrorZoneNotFound)
		shared.LogError("error finding zone", LogService, "AddTablesToZone", err, zoneID)
		return err
	}

	return s.zones.AddTables(zone, tables)
}

func (s service) RemoveTables(zoneID string, tables []uint) error {
	zone, err := s.zones.GetZone(zoneID)
	if err != nil {
		return err
	}

	if zone == nil {
		err := fmt.Errorf(ErrorZoneNotFound)
		shared.LogError("error finding zone", LogService, "RemoveTablesFromZone", err, zoneID)
		return err
	}

	return s.zones.RemoveTables(zone, tables)
}

func (s service) EnableZone(zoneID string) (*Zone, error) {
	return s.zones.Enable(zoneID)
}

var _ Service = &service{}
