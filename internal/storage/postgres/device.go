package postgres

import (
	"github.com/subzerobo/ratatoskr/internal/services/devices"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// First  string `gorm:"index:idx_name,unique"`
// Second string `gorm:"index:idx_name,unique"

type device struct {
	ID                uint   `gorm:"primary_key"`
	UUID              string `gorm:"type:uuid;not null;default:uuid_generate_v4()"`
	Identifier        string `gorm:"size:512;uniqueIndex:idx_adid_identifier;uniqueIndex:idx_identifier"` // PushToken, Email, PhoneNumber
	DeviceType        string `gorm:"size:50"`                                                             // Android / iOS / Web
	Language          string `gorm:"size:7"`                                                              // Two-Letters except for Chinese
	Timezone          int
	AppVersion        string `gorm:"size:32"`
	DeviceVendor      string `gorm:"size:255"`
	DeviceModel       string `gorm:"size:255"`
	DeviceOS          string `gorm:"size:255"`
	DeviceOSVersion   string `gorm:"size:255"`
	ADID              string `gorm:"size:255;uniqueIndex:idx_adid_identifier"`
	SDK               string `gorm:"size:255"`
	SessionCount      int
	NotificationTypes int
	Long              float32
	Lat               float32
	Country           string    `gorm:"size:2"`
	ExternalUserID    string    `gorm:"size:255;index"`
	CreatedAt         time.Time `gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp"`
	ApplicationID     uint      `gorm:"index"`
	Application       application
	Tags              []tag `gorm:"foreignKey:DeviceID"`
	BadgeCount        int
	AmountSpent       float32
}

type tag struct {
	DeviceID uint   `gorm:"index;uniqueIndex:idx_mix"`
	Key      string `gorm:"size:255;index;uniqueIndex:idx_mix"`
	Value    string `gorm:"size:255;index"`
	Device   device
}

func (d device) ToServiceModel() *devices.DeviceModel {
	dm := devices.DeviceModel{
		ID:                d.ID,
		UUID:              &d.UUID,
		DeviceType:        &d.DeviceType,
		Identifier:        &d.Identifier,
		Language:          &d.Language,
		Timezone:          d.Timezone,
		AppVersion:        &d.AppVersion,
		DeviceVendor:      &d.DeviceVendor,
		DeviceModel:       &d.DeviceModel,
		DeviceOS:          &d.DeviceOS,
		DeviceOSVersion:   &d.DeviceOSVersion,
		ADID:              &d.ADID,
		SDK:               &d.SDK,
		SessionCount:      &d.SessionCount,
		NotificationTypes: &d.NotificationTypes,
		Long:              &d.Long,
		Lat:               &d.Lat,
		Country:           &d.Country,
		ExternalUserID:    &d.ExternalUserID,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
		ApplicationID:     &d.ApplicationID,
		BadgeCount:        &d.BadgeCount,
		AmountSpent:       &d.AmountSpent,
	}
	dm.Tags = make(map[string]string)
	for _, t := range d.Tags {
		dm.Tags[t.Key] = t.Value
	}
	return &dm
}

func (r *repository) UpsertDevice(model devices.DeviceModel) (*devices.DeviceModel, error) {
	dev := device{
		DeviceType:        *model.DeviceType,
		Identifier:        *model.Identifier,
		Language:          *model.Language,
		Timezone:          model.Timezone,
		AppVersion:        *model.AppVersion,
		DeviceVendor:      *model.DeviceVendor,
		DeviceModel:       *model.DeviceModel,
		DeviceOS:          *model.DeviceOS,
		DeviceOSVersion:   *model.DeviceOSVersion,
		ADID:              *model.ADID,
		SDK:               *model.SDK,
		SessionCount:      *model.SessionCount,
		NotificationTypes: *model.NotificationTypes,
		Long:              *model.Long,
		Lat:               *model.Lat,
		Country:           *model.Country,
		ExternalUserID:    *model.ExternalUserID,
		ApplicationID:     *model.ApplicationID,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "identifier"}, {Name: "ad_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"language":           dev.Language,
				"timezone":           dev.Timezone,
				"app_version":        dev.AppVersion,
				"device_vendor":      dev.DeviceVendor,
				"device_model":       dev.DeviceModel,
				"device_os":          dev.DeviceOS,
				"device_os_version":  dev.DeviceOSVersion,
				"sdk":                dev.SDK,
				"session_count":      dev.SessionCount,
				"notification_types": dev.NotificationTypes,
				"long":               dev.Long,
				"lat":                dev.Lat,
				"country":            dev.Country,
				"external_user_id":   dev.ExternalUserID,
			}),
		}).Create(&dev).Error
		if err != nil {
			return errors.Wrapf(err, "failed to insert device record %v", dev)
		}

		for k, v := range model.Tags {
			tagM := tag{
				DeviceID: dev.ID,
				Key:      k,
				Value:    v,
			}
			if v == "" {
				// Remove tag if value is empty
				err = tx.Where("device_id = ? AND key = ?", tagM.DeviceID, tagM.Key).Delete(tag{}).Error
			} else {
				// Upsert tag if value is not empty
				err = tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "device_id"}, {Name: "key"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"value": tagM.Value,
					}),
				}).Create(&tagM).Error
			}
			if err != nil {
				return errors.Wrapf(err, "failed to insert tag record %v", tagM)
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to commit changes")
	}

	return dev.ToServiceModel(), nil
}

func (r *repository) GetDevice(uuid string, applicationID uint) (*devices.DeviceModel, error) {
	var item device
	err := r.db.Where("uuid = ? AND application_id = ?", uuid, applicationID).First(&item).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}

	return item.ToServiceModel(), nil
}

func (r *repository) GetDevices(applicationID uint, lastID uint, limit int) ([]*devices.DeviceModel, error) {
	var items []device
	err := r.db.Where("application_id = ? AND id > ? ", applicationID, lastID).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	var result []*devices.DeviceModel
	for _, item := range items {
		result = append(result, item.ToServiceModel())
	}
	return result, nil
}

func (r *repository) UpdatePartial(model devices.DeviceModel) (*devices.DeviceModel, error) {
	dev := device{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("uuid = ?", model.UUID).First(&dev).Error
		if err != nil {
			return getProcessedDBError(err)
		}

		dev.Timezone = model.Timezone
		if model.DeviceModel != nil {
			dev.DeviceModel = *model.DeviceModel
		}
		if model.Identifier != nil {
			dev.Identifier = *model.Identifier
		}
		if model.Language != nil {
			dev.Language = *model.Language
		}
		if model.AppVersion != nil {
			dev.AppVersion = *model.AppVersion
		}
		if model.DeviceVendor != nil {
			dev.DeviceVendor = *model.DeviceVendor
		}
		if model.DeviceModel != nil {
			dev.DeviceModel = *model.DeviceModel
		}
		if model.DeviceOS != nil {
			dev.DeviceOS = *model.DeviceOS
		}
		if model.DeviceOSVersion != nil {
			dev.DeviceOSVersion = *model.DeviceOSVersion
		}
		if model.ADID != nil {
			dev.ADID = *model.ADID
		}
		if model.SDK != nil {
			dev.SDK = *model.SDK
		}
		if model.SessionCount != nil {
			dev.SessionCount = *model.SessionCount
		}
		if model.NotificationTypes != nil {
			dev.NotificationTypes = *model.NotificationTypes
		}
		if model.SessionCount != nil {
			dev.SessionCount = *model.SessionCount
		}
		if model.Long != nil {
			dev.Long = *model.Long
		}
		if model.Lat != nil {
			dev.Lat = *model.Lat
		}
		if model.Country != nil {
			dev.Country = *model.Country
		}
		if model.ExternalUserID != nil {
			dev.ExternalUserID = *model.ExternalUserID
		}
		if model.BadgeCount != nil {
			dev.BadgeCount = *model.BadgeCount
		}
		if model.AmountSpent != nil {
			dev.AmountSpent = *model.AmountSpent
		}

		for k, v := range model.Tags {
			tagM := tag{
				DeviceID: dev.ID,
				Key:      k,
				Value:    v,
			}
			if v == "" {
				// Remove tag if value is empty
				err = tx.Where("device_id = ? AND key = ?", tagM.DeviceID, tagM.Key).Delete(tag{}).Error
			} else {
				// Upsert tag if value is not empty
				err = tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "device_id"}, {Name: "key"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"value": tagM.Value,
					}),
				}).Create(&tagM).Error
			}
			if err != nil {
				return errors.Wrapf(err, "failed to insert tag record %v", tagM)
			}
		}

		err = tx.Save(dev).Error
		return err
	})

	if err != nil {
		return nil, err
	}
	return dev.ToServiceModel(), nil
}

func (r *repository) UpdateDeviceTagsByUser(applicationID uint, externalUserID string, Tags map[string]string) error {
	var items []device
	err := r.db.Where("application_id = ? AND external_user_id = ? ", applicationID, externalUserID).Find(&items).Error
	if err != nil {
		return getProcessedDBError(err)
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			for k, v := range Tags {
				tagM := tag{
					DeviceID: item.ID,
					Key:      k,
					Value:    v,
				}
				if v == "" {
					// Remove tag if value is empty
					err = tx.Where("device_id = ? AND key = ?", tagM.DeviceID, tagM.Key).Delete(tag{}).Error
				} else {
					// Upsert tag if value is not empty
					err = tx.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "device_id"}, {Name: "key"}},
						DoUpdates: clause.Assignments(map[string]interface{}{
							"value": tagM.Value,
						}),
					}).Create(&tagM).Error
				}
				if err != nil {
					return errors.Wrapf(err, "failed to insert tag record %v", tagM)
				}
			}
		}
		return nil
	})
	return err

	return nil
}
