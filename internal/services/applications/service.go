package applications

import (
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/utils"
)

type Service interface {
	Create(model ApplicationModel) (*ApplicationModel, error)
	List(accountID uint) ([]*ApplicationModel, error)
	ListAll() ([]*ApplicationModel, error)
	UpdateAuthKey(accountID uint, UUID string) (string, error)
	UpdateIdentityVerification(accountID uint, UUID string, status bool) error
	Details(accountID uint, UUID string) (*ApplicationModel, error)
	Delete(accountID uint, UUID string) error
}

type service struct {
	repository Repository
}

func CreateService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(model ApplicationModel) (*ApplicationModel, error) {
	model.AuthKey = utils.RandomString(64)
	res, err := s.repository.CreateApplication(model)
	if err != nil {
		return nil, err
	}
	
	return res, nil
	
}

func (s service) List(accountID uint) ([]*ApplicationModel, error) {
	res, err := s.repository.GetApplicationsByAccountID(accountID)
	if err != nil {
		if errors.HasKind(err, errors.NotFound) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (s service) ListAll() ([]*ApplicationModel, error) {
	res, err := s.repository.GetAllApplications()
	if err != nil {
		if errors.HasKind(err, errors.NotFound) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (s service) Details(accountID uint, UUID string) (*ApplicationModel, error) {
	res, err := s.repository.GetAccountApplicationByUUID(accountID,UUID)
	return res, err
}

func (s service) UpdateAuthKey(accountID uint, UUID string) (string, error) {
	newToken := utils.RandomString(64)
	err := s.repository.UpdateAuthKey(accountID,UUID, newToken)
	if err != nil {
		return "", err
	}
	return newToken,nil
}

func (s service) UpdateIdentityVerification(accountID uint, UUID string, status bool) error {
	return s.repository.UpdateIdentityVerification(accountID,UUID, status)
}

func (s service) Delete(accountID uint, UUID string) error {
	panic("implement me")
}
