package postgres

import (
	applications2 "github.com/subzerobo/ratatoskr/internal/services/applications"
	devices2 "github.com/subzerobo/ratatoskr/internal/services/devices"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"time"
)

type application struct {
	ID                   uint      `gorm:"primary_key"`
	UUID                 string    `gorm:"type:uuid; not null;default:uuid_generate_v4()"`
	Name                 string    `gorm:"size:255"`
	FCMSenderID          string    `gorm:"uniqueIndex;size:255"`
	FCMAdminJSON         string    `gorm:"size:5000"`
	URL                  string    `gorm:"size:255"`
	AuthKey              string    `gorm:"uniqueIndex;size:64"`
	IdentityVerification bool      `gorm:"default:false"`
	CreatedAt            time.Time `gorm:"default:current_timestamp"`
	UpdatedAt            time.Time `gorm:"default:current_timestamp"`
	AccountID            uint      `gorm:"index"`
	Account              account
}

func (a application) ToServiceModel() *applications2.ApplicationModel {
	return &applications2.ApplicationModel{
		ID:                   a.ID,
		UUID:                 a.UUID,
		Name:                 a.Name,
		FCMSenderID:          a.FCMSenderID,
		FCMAdminJSON:         a.FCMAdminJSON,
		URL:                  a.URL,
		AuthKey:              a.AuthKey,
		IdentityVerification: a.IdentityVerification,
		CreatedAt:            a.CreatedAt,
		UpdatedAt:            a.UpdatedAt,
		AccountID:            0,
	}
}

func (r *repository) CreateApplication(model applications2.ApplicationModel) (*applications2.ApplicationModel, error) {
	app := application{
		Name:         model.Name,
		FCMSenderID:  model.FCMSenderID,
		FCMAdminJSON: model.FCMAdminJSON,
		URL:          model.URL,
		AccountID:    model.AccountID,
	}
	
	err := r.db.Create(&app).Error
	if err != nil {
		return nil, errors.WithKindCtx(err, "failed to insert record to database", errors.InternalServerError, nil)
	}
	return app.ToServiceModel(), nil
}

func (r *repository) GetApplicationsByAccountID(accountID uint) ([]*applications2.ApplicationModel, error) {
	var items []application
	err := r.db.Where("account_id = ?", accountID).Find(&items).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	var result []*applications2.ApplicationModel
	for _, item := range items {
		result = append(result, item.ToServiceModel())
	}
	return result, nil
}

func (r *repository) GetAllApplications() ([]*applications2.ApplicationModel, error) {
	var items []application
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	var result []*applications2.ApplicationModel
	for _, item := range items {
		result = append(result, item.ToServiceModel())
	}
	return result, nil
}

func (r *repository) GetAccountApplicationByUUID(accountID uint, UUID string) (*applications2.ApplicationModel, error) {
	var item application
	err := r.db.Where("uuid = ? AND account_id = ?", UUID, accountID).First(&item).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	
	return item.ToServiceModel(), nil
}

func (r *repository) UpdateAuthKey(accountID uint, UUID string, AuthKey string) error {
	return r.db.Model(&application{}).
		Where("account_id = ? AND uuid = ?", accountID, UUID).
		Update("auth_key", AuthKey).Error
}

func (r *repository) UpdateIdentityVerification(accountID uint, UUID string, status bool) error {
	return r.db.Model(&application{}).
		Where("account_id = ? AND uuid = ?", accountID, UUID).
		Update("identity_verification", status).Error
}

func (r *repository) GetApplicationByUUID(UUID string) (*devices2.DeviceApplicationModel, error) {
	var item application
	err := r.db.Where("uuid = ?", UUID).First(&item).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	
	return &devices2.DeviceApplicationModel{
		ID:                   item.ID,
		UUID:                 item.UUID,
		AuthKey:              item.AuthKey,
		IdentityVerification: item.IdentityVerification,
		AccountID:            item.AccountID,
	}, nil
}