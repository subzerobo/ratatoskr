package authentication

type Repository interface {
	CreateAccount(model AccountModel) (*AccountModel, error)
	GetAccountByEmail(email string) (*AccountModel, error)
}

type StateStore interface {
	SetState(key, value string) error
	GetState(key string) (string, error)
}
