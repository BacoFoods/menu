package channel

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository string = "pkg/channel/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new channel in database
func (r *DBRepository) Create(channel *Channel) (*Channel, error) {
	if err := r.db.Save(channel).Error; err != nil {
		shared.LogError("error creating channel", LogDBRepository, "Create", err, channel)
		return nil, err
	}
	return channel, nil
}

// Find method for find channels in database
func (r *DBRepository) Find(filters map[string]string) ([]Channel, error) {
	var channels []Channel
	if err := r.db.Find(&channels, filters).Error; err != nil {
		shared.LogError("error getting channels", LogDBRepository, "Find", err, filters)
		return nil, err
	}

	return channels, nil
}

// Get method for get a channel in database
func (r *DBRepository) Get(channelID string) (*Channel, error) {
	if channelID == "" {
		shared.LogWarn("error getting channel", LogDBRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var channel Channel
	if err := r.db.First(&channel, channelID).Error; err != nil {
		shared.LogError("error getting channel", LogDBRepository, "Get", err, channelID)
		return nil, err
	}
	return &channel, nil
}

// Update method for update a channel in database
func (r *DBRepository) Update(channel *Channel) (*Channel, error) {
	var channelDB Channel
	if err := r.db.First(&channelDB, channel.ID).Error; err != nil {
		shared.LogError("error getting channel", LogDBRepository, "Update", err, channel)
		return nil, err
	}
	if err := r.db.Model(&channelDB).Updates(channel).Error; err != nil {
		shared.LogError("error updating channel", LogDBRepository, "Update", err, channel)
		return nil, err
	}
	return &channelDB, nil
}

// Delete method for delete a channel in database
func (r *DBRepository) Delete(channelID string) (*Channel, error) {
	var channel Channel

	if err := r.db.First(&channel, channelID).Error; err != nil {
		shared.LogError("error getting channel", LogDBRepository, "Delete", err, channelID)
		return nil, err
	}

	if err := r.db.Delete(&channel).Error; err != nil {
		shared.LogError("error deleting channel", LogDBRepository, "Delete", err, channel)
		return nil, err
	}

	return &channel, nil
}

// FindByIDs method for find channels by ids in database
func (r *DBRepository) FindByIDs(channelIDs []string) ([]Channel, error) {
	var channels []Channel

	if err := r.db.Find(&channels, channelIDs).Error; err != nil {
		shared.LogError("error getting channels", LogDBRepository, "FindByIDs", err, channelIDs)
		return nil, err
	}

	return channels, nil
}
