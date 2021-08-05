package postgres

import (
	"github.com/subzerobo/ratatoskr/internal/services/applications"
	"github.com/subzerobo/ratatoskr/internal/services/devices"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"gorm.io/gorm"
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

func (a application) ToServiceModel() *applications.ApplicationModel {
	return &applications.ApplicationModel{
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

type androidGroup struct {
	ID            uint                   `gorm:"primary_key"`
	GroupName     string                 `gorm:"size:255"`
	GroupUUID     string                 `gorm:"type:uuid; not null;default:uuid_generate_v4()"`
	CreatedAt     time.Time              `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time              `gorm:"default:current_timestamp"`
	Categories    []androidGroupCategory `gorm:"foreignKey:AndroidGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ApplicationID uint                   `gorm:"index"`
	Application   application
}

func (a androidGroup) ToServiceModel() *applications.AndroidGroupModel {
	res := &applications.AndroidGroupModel{
		ID:            a.ID,
		ApplicationID: a.ApplicationID,
		GroupName:     a.GroupName,
		GroupUUID:     a.GroupUUID,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
	for _, cat := range a.Categories {
		res.Categories = append(res.Categories, *cat.ToServiceModel())
	}
	return res
}

type androidGroupCategory struct {
	ID                  uint      `gorm:"primary_key"`
	CategoryUUID        string    `gorm:"type:uuid; not null;default:uuid_generate_v4()"`
	CategoryName        string    `gorm:"size:255"`
	CategoryDescription string    `gorm:"size:255"`
	Priority            string    `gorm:"size:1;default:3"`
	Sound               int       `gorm:"default:1"`
	SoundName           string    `gorm:"size:255"`
	Vibration           int       `gorm:"default:1"`
	VibrationPattern    string    `gorm:"size:255"`
	Led                 int       `gorm:"default:1"`
	LedColor            string    `gorm:"size:8"`
	EnableBadge         int       `gorm:"default:1"`
	LockScreen          int       `gorm:"default:0"`
	CreatedAt           time.Time `gorm:"default:current_timestamp"`
	UpdatedAt           time.Time `gorm:"default:current_timestamp"`
	AndroidGroupID      uint      `gorm:"index"`
	AndroidGroup        androidGroup
}

func (a androidGroupCategory) ToServiceModel() *applications.AndroidGroupCategoryModel {
	return &applications.AndroidGroupCategoryModel{
		ID:                  a.ID,
		CategoryUUID:        a.CategoryUUID,
		CategoryName:        a.CategoryName,
		CategoryDescription: a.CategoryDescription,
		Priority:            a.Priority,
		Sound:               a.Sound,
		SoundName:           a.SoundName,
		Vibration:           a.Vibration,
		VibrationPattern:    a.VibrationPattern,
		Led:                 a.Led,
		LedColor:            a.LedColor,
		EnableBadge:         a.EnableBadge,
		LockScreen:          a.LockScreen,
		CreatedAt:           a.CreatedAt,
		UpdatedAt:           a.UpdatedAt,
		AndroidGroupID:      a.AndroidGroupID,
	}
}

func (r *repository) CreateApplication(model applications.ApplicationModel) (*applications.ApplicationModel, error) {
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

func (r *repository) GetApplicationsByAccountID(accountID uint) ([]*applications.ApplicationModel, error) {
	var items []application
	err := r.db.Where("account_id = ?", accountID).Find(&items).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	var result []*applications.ApplicationModel
	for _, item := range items {
		result = append(result, item.ToServiceModel())
	}
	return result, nil
}

func (r *repository) GetAllApplications() ([]*applications.ApplicationModel, error) {
	var items []application
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	var result []*applications.ApplicationModel
	for _, item := range items {
		result = append(result, item.ToServiceModel())
	}
	return result, nil
}

func (r *repository) GetAccountApplicationByUUID(accountID uint, UUID string) (*applications.ApplicationModel, error) {
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

func (r *repository) GetApplicationByUUID(UUID string) (*devices.DeviceApplicationModel, error) {
	var item application
	err := r.db.Where("uuid = ?", UUID).First(&item).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}

	return &devices.DeviceApplicationModel{
		ID:                   item.ID,
		UUID:                 item.UUID,
		AuthKey:              item.AuthKey,
		IdentityVerification: item.IdentityVerification,
		AccountID:            item.AccountID,
	}, nil
}

func (r *repository) GetApplicationModelByUUID(UUID string) (*applications.ApplicationModel, error) {
	var item application
	err := r.db.Where("uuid = ?", UUID).First(&item).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}

	return item.ToServiceModel(), nil
}

func (r *repository) GetAndroidGroups(appId uint) ([]*applications.AndroidGroupModel, error) {
	var items []androidGroup
	err := r.db.Preload("Categories").Where("application_id = ?", appId).Find(&items).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	var result []*applications.AndroidGroupModel
	for _, item := range items {
		result = append(result, item.ToServiceModel())
	}
	return result, nil
}

func (r *repository) CreateAndroidGroup(model applications.AndroidGroupModel) (*applications.AndroidGroupModel, error) {
	agp := androidGroup{
		GroupName:     model.GroupName,
		ApplicationID: model.ApplicationID,
	}

	err := r.db.Create(&agp).Error
	if err != nil {
		return nil, getProcessedDBError(err)
	}
	return agp.ToServiceModel(), nil
}

func (r *repository) UpdateAndroidGroup(model applications.AndroidGroupModel) (*applications.AndroidGroupModel, error) {
	var agp androidGroup
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("group_uuid = ? and application_id = ?", model.GroupUUID, model.ApplicationID).First(&agp).Error
		if err != nil {
			return getProcessedDBError(err)
		}

		agp.GroupName = model.GroupName
		err = tx.Save(agp).Error
		return err
	})
	return agp.ToServiceModel(), err
}

func (r *repository) DeleteAndroidGroup(appId uint, uuid string) error {
	return r.db.Where("group_uuid = ? and application_id = ?", uuid, appId).Delete(androidGroup{}).Error
}

func (r *repository) CreateAndroidCategory(appID uint, groupUUID string, model applications.AndroidGroupCategoryModel) error {
	// Find the Group first
	var agp androidGroup
	err := r.db.Where("group_uuid = ? and application_id = ?", groupUUID, appID).First(&agp).Error
	if err != nil {
		return getProcessedDBError(err)
	}

	agpc := androidGroupCategory{
		CategoryName:        model.CategoryName,
		CategoryDescription: model.CategoryDescription,
		Priority:            model.Priority,
		Sound:               model.Sound,
		SoundName:           model.SoundName,
		Vibration:           model.Vibration,
		VibrationPattern:    model.VibrationPattern,
		Led:                 model.Led,
		LedColor:            model.LedColor,
		EnableBadge:         model.EnableBadge,
		LockScreen:          model.LockScreen,
		AndroidGroupID:      agp.ID,
	}

	err = r.db.Create(&agpc).Error
	return getProcessedDBError(err)
}

func (r *repository) UpdateAndroidCategory(appID uint, groupUUID string, model applications.AndroidGroupCategoryModel) error {
	// Find the Group first
	var agp androidGroup
	err := r.db.Where("group_uuid = ? and application_id = ?", groupUUID, appID).First(&agp).Error
	if err != nil {
		return getProcessedDBError(err)
	}

	var agpc androidGroupCategory
	err = r.db.Where("category_uuid = ? and android_group_id = ?", model.CategoryUUID, agp.ID).First(&agpc).Error
	if err != nil {
		return getProcessedDBError(err)
	}

	agpc.CategoryName = model.CategoryName
	agpc.CategoryDescription = model.CategoryDescription
	agpc.Led = model.Led
	agpc.LedColor = model.LedColor
	agpc.Vibration = model.Vibration
	agpc.VibrationPattern = model.VibrationPattern
	agpc.Sound = model.Sound
	agpc.SoundName = model.SoundName
	agpc.Priority = model.Priority
	agpc.LockScreen = model.LockScreen
	agpc.EnableBadge = model.EnableBadge

	return r.db.Save(&agpc).Error
}

func (r *repository) DeleteAndroidCategory(appID uint, groupUUID string, categoryUUID string, ) error {
	// Find the Group first
	var agp androidGroup
	err := r.db.Where("group_uuid = ? and application_id = ?", groupUUID, appID).First(&agp).Error
	if err != nil {
		return getProcessedDBError(err)
	}

	return r.db.Where("category_uuid = ? and android_group_id = ?", categoryUUID, agp.ID).Delete(androidGroupCategory{}).Error
}
