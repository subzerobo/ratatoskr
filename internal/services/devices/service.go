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
	Update(model DeviceModel) (*DeviceModel, error)
	UpdateUserTags(AppUUID string, externalUserID string, tags map[string]string) error
	Get(UUID string, AppUUID string) (*DeviceModel, error)
	GetList(AppUUID string, paging utils.MorePaging) ([]*DeviceModel, error)
}

type service struct {
	repository Repository
}

func CreateService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Upsert(model DeviceModel, AppUUID string) (*DeviceModel, error) {
	// Check if APP_ID is valid
	app, err := s.repository.GetApplicationByUUID(AppUUID)
	if err != nil {
		return nil, err
	}

	// Check if Identity Check is enabled
	if app.IdentityVerification && model.ExternalUserID != nil {
		if !utils.CheckHMACHash(*model.ExternalUserID, *model.ExternalUserIDHash, app.AuthKey) {
			return nil, errors.WithKindCtx(ErrInvalidHash, "", errors.Unauthorized, nil)
		}
	}

	// Update Application ID
	model.ApplicationID = &app.ID

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

func (s service) GetList(AppUUID string, paging utils.MorePaging) ([]*DeviceModel, error) {
	app, err := s.repository.GetApplicationByUUID(AppUUID)
	if err != nil {
		return nil, err
	}
	return s.repository.GetDevices(app.ID, paging.LastID, paging.Size)
}

func (s service) Update(model DeviceModel) (*DeviceModel, error) {
	res, err := s.repository.UpdatePartial(model)
	return res, err
}

func (s service) UpdateUserTags(AppUUID string, externalUserID string, tags map[string]string) error {
	app, err := s.repository.GetApplicationByUUID(AppUUID)
	if err != nil {
		return  err
	}
	return s.repository.UpdateDeviceTagsByUser(app.ID, externalUserID, tags)
}