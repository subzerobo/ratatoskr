package postgres

import (
	"database/sql"
	authentication2 "github.com/subzerobo/ratatoskr/internal/services/authentication"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"time"
)

type account struct {
	ID                 uint   `gorm:"primary_key"`
	UUID               string `gorm:"type:uuid; not null;default:uuid_generate_v4()"`
	Email              string `gorm:"uniqueIndex;size:50"`
	EncryptedPassword  string `gorm:"size:255"`
	OAuthProvider      string `gorm:"size:255"`
	OAuthUID           string `gorm:"size:255"`
	Picture            string `gorm:"size:255"`
	CompanyName        string `gorm:"size:255"`
	IsSuperUser        bool   `gorm:"not null;default:false"`
	LastLoginDate      sql.NullTime
	Active             bool `gorm:"not null;default:true"`
	Confirmed          bool `gorm:"not null;default:false"`
	ConfirmationToken  string
	ConfirmationSentAt sql.NullTime
	CreatedAt          time.Time `gorm:"default:current_timestamp"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp"`
}

func (a account) ToServiceModel() *authentication2.AccountModel {
	return &authentication2.AccountModel{
		ID:                 a.ID,
		UUID:               a.UUID,
		Email:              a.Email,
		EncryptedPassword:  a.EncryptedPassword,
		OAuthProvider:      a.OAuthProvider,
		OAuthUID:           a.OAuthUID,
		Picture:            a.Picture,
		CompanyName:        a.CompanyName,
		IsSuperUser:        a.IsSuperUser,
		LastLoginDate:      a.LastLoginDate,
		Active:             a.Active,
		Confirmed:          a.Confirmed,
		ConfirmationToken:  a.ConfirmationToken,
		ConfirmationSentAt: a.ConfirmationSentAt,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

func (r *repository) CreateAccount(model authentication2.AccountModel) (*authentication2.AccountModel, error) {
	acc := account{
		Email:              model.Email,
		EncryptedPassword:  model.EncryptedPassword,
		OAuthProvider:      model.OAuthProvider,
		OAuthUID:           model.OAuthUID,
		Picture:            model.Picture,
		CompanyName:        model.CompanyName,
		IsSuperUser:        model.IsSuperUser,
		LastLoginDate:      sql.NullTime{},
		Active:             model.Active,
		Confirmed:          model.Confirmed,
		ConfirmationToken:  model.ConfirmationToken,
		ConfirmationSentAt: model.ConfirmationSentAt,
	}
	err := r.db.Create(&acc).Error
	if err != nil {
		return nil, errors.WithKindCtx(err, "failed to insert record to database", errors.InternalServerError, nil)
	}
	return acc.ToServiceModel(), nil
}

func (r *repository) GetAccountByEmail(email string) (*authentication2.AccountModel, error) {
	var acc account
	err := r.db.Where("email = ?", email).First(&acc).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	return acc.ToServiceModel(), nil
}
