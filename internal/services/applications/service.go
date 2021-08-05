package applications

import (
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/utils"
)

var (
	ErrInvalidApplicationAuthKey = errors.New("application uuid/auth token is invalid")
)

type Service interface {
	Create(model ApplicationModel) (*ApplicationModel, error)
	List(accountID uint) ([]*ApplicationModel, error)
	ListAll() ([]*ApplicationModel, error)
	UpdateAuthKey(accountID uint, UUID string) (string, error)
	UpdateIdentityVerification(accountID uint, UUID string, status bool) error
	Details(accountID uint, UUID string) (*ApplicationModel, error)
	Delete(accountID uint, UUID string) error
	CheckApplicationToken(authKey string, UUID string) error

	GetAndroidGroups(UUID string) ([]*AndroidGroupModel, error)
	CreateAndroidGroup(accountID uint, aUUID string, Name string) error
	UpdateAndroidGroup(accountID uint, aUUID string, gUUID string, Name string) error
	DeleteAndroidGroup(accountID uint, aUUID string, gUUID string) error

	CreateAndroidCategory(accountID uint, aUUID string, gUUID string, model AndroidGroupCategoryModel) error
	UpdateAndroidCategory(accountID uint, aUUID string, gUUID string, model AndroidGroupCategoryModel) error
	DeleteAndroidCategory(accountID uint, aUUID string, gUUID string, cUUID string) error
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
	res, err := s.repository.GetAccountApplicationByUUID(accountID, UUID)
	return res, err
}

func (s service) UpdateAuthKey(accountID uint, UUID string) (string, error) {
	newToken := utils.RandomString(64)
	err := s.repository.UpdateAuthKey(accountID, UUID, newToken)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

func (s service) UpdateIdentityVerification(accountID uint, UUID string, status bool) error {
	return s.repository.UpdateIdentityVerification(accountID, UUID, status)
}

func (s service) Delete(accountID uint, UUID string) error {
	panic("implement me")
}

func (s service) CheckApplicationToken(authKey string, UUID string) error {
	res, err := s.repository.GetApplicationModelByUUID(UUID)
	if err != nil {
		return err
	}
	if res.AuthKey != authKey {
		return errors.WithKindCtx(ErrInvalidApplicationAuthKey, "", errors.Unauthorized, nil)
	}
	return nil
}

func (s service) GetAndroidGroups(UUID string) ([]*AndroidGroupModel, error) {
	res, err := s.repository.GetApplicationModelByUUID(UUID)
	if err != nil {
		return nil, err
	}
	groups, err := s.repository.GetAndroidGroups(res.ID)
	return groups, err

}

func (s service) CreateAndroidGroup(accountID uint, aUUID string, Name string) error {
	// Get Application By AccountID and UUID
	res, err := s.repository.GetAccountApplicationByUUID(accountID, aUUID)
	if err != nil {
		return err
	}

	_, err = s.repository.CreateAndroidGroup(AndroidGroupModel{
		ApplicationID: res.ID,
		GroupName:     Name,
	})

	return err
}

func (s service) UpdateAndroidGroup(accountID uint, aUUID string, gUUID string, Name string) error {
	// Get Application By AccountID and UUID
	res, err := s.repository.GetAccountApplicationByUUID(accountID, aUUID)
	if err != nil {
		return err
	}

	_, err = s.repository.UpdateAndroidGroup(AndroidGroupModel{
		ApplicationID: res.ID,
		GroupUUID:     gUUID,
		GroupName:     Name,
	})
	return err
}

func (s service) DeleteAndroidGroup(accountID uint, aUUID string, gUUID string) error {
	// Get Application By AccountID and UUID
	res, err := s.repository.GetAccountApplicationByUUID(accountID, aUUID)
	if err != nil {
		return err
	}
	return s.repository.DeleteAndroidGroup(res.ID, gUUID)
}

func (s service) CreateAndroidCategory(accountID uint, aUUID string, gUUID string, model AndroidGroupCategoryModel) error {
	res, err := s.repository.GetAccountApplicationByUUID(accountID, aUUID)
	if err != nil {
		return err
	}

	err = s.repository.CreateAndroidCategory(res.ID, gUUID,model)
	return err
}

func (s service) UpdateAndroidCategory(accountID uint, aUUID string, gUUID string, model AndroidGroupCategoryModel) error {
	res, err := s.repository.GetAccountApplicationByUUID(accountID, aUUID)
	if err != nil {
		return err
	}

	err = s.repository.UpdateAndroidCategory(res.ID, gUUID, model)
	return err
}

func (s service) DeleteAndroidCategory(accountID uint, aUUID string, gUUID string, cUUID string) error {
	res, err := s.repository.GetAccountApplicationByUUID(accountID, aUUID)
	if err != nil {
		return err
	}

	err = s.repository.DeleteAndroidCategory(res.ID, gUUID, cUUID)
	return err
}