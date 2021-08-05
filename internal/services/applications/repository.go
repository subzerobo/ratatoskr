package applications

type Repository interface {
	CreateApplication(model ApplicationModel) (*ApplicationModel, error)
	GetApplicationsByAccountID(accountID uint) ([]*ApplicationModel, error)
	GetAllApplications() ([]*ApplicationModel, error)
	GetAccountApplicationByUUID(accountID uint, UUID string) (*ApplicationModel, error)
	GetApplicationModelByUUID(UUID string) (*ApplicationModel, error)
	UpdateAuthKey(accountID uint, UUID string, AuthKey string) error
	UpdateIdentityVerification(accountID uint, UUID string, status bool) error

	GetAndroidGroups(appId uint) ([]*AndroidGroupModel, error)
	CreateAndroidGroup(model AndroidGroupModel) (*AndroidGroupModel, error)
	UpdateAndroidGroup(model AndroidGroupModel) (*AndroidGroupModel, error)
	DeleteAndroidGroup(appId uint, uuid string) error

	CreateAndroidCategory(appID uint, groupUUID string, model AndroidGroupCategoryModel) error
	UpdateAndroidCategory(appID uint, groupUUID string, model AndroidGroupCategoryModel) error
	DeleteAndroidCategory(appID uint, groupUUID string, categoryUUID string,) error

}
