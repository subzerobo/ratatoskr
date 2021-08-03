package devices

import (
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/utils"
)

var (
	ErrInvalidHash = errors.New("invalid external user id hash")
)

type Service interface {
	Upsert(model DeviceModel, AppUUID string) (*DeviceModel, error)
	Get(UUID string, AppUUID string) (*DeviceModel, error)
}

type service struct {
	repository Repository
}

func CreateService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Upsert(model DeviceModel,AppUUID string) (*DeviceModel, error) {
	// Check if APP_ID is valid
	app, err := s.repository.GetApplicationByUUID(AppUUID)
	if err != nil {
		return nil, err
	}
	
	// Check if Identity Check is enabled
	if app.IdentityVerification && model.ExternalUserID != "" {
		if !utils.CheckHMACHash(model.ExternalUserID, model.ExternalUserIDHash, app.AuthKey) {
			return nil, errors.WithKindCtx(ErrInvalidHash, "", errors.Unauthorized, nil)
		}
	}
	
	// Update Application ID
	model.ApplicationID = app.ID
	
	// Add device record
	res, err := s.repository.UpsertDevice(model)
	return res, err
}

func (s service) Get(UUID string, AppUUID string) (*DeviceModel, error) {
	app, err := s.repository.GetApplicationByUUID(AppUUID)
	if err != nil {
		return nil, err
	}
	return s.repository.GetDevice(UUID, app.ID)
}
