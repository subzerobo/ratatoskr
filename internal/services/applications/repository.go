package applications

type Repository interface {
	CreateApplication(model ApplicationModel) (*ApplicationModel, error)
	GetApplicationsByAccountID(accountID uint) ([]*ApplicationModel, error)
	GetAllApplications() ([]*ApplicationModel, error)
	GetAccountApplicationByUUID(accountID uint, UUID string) (*ApplicationModel, error)
	UpdateAuthKey(accountID uint, UUID string, AuthKey string) error
	UpdateIdentityVerification(accountID uint, UUID string, status bool) error
	// GetApplicationByID(id uint) (*ApplicationModel, error)
}
