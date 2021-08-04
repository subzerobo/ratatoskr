package devices

type Repository interface {
	UpsertDevice(model DeviceModel) (*DeviceModel, error)
	UpdatePartial(model DeviceModel) (*DeviceModel, error)
	GetDevice(uuid string, applicationID uint) (*DeviceModel, error)
	GetDevices(applicationID uint, lastID uint, limit int) ([]*DeviceModel, error)
	GetApplicationByUUID(uuid string) (*DeviceApplicationModel, error)
	UpdateDeviceTagsByUser(applicationID uint, externalUserID string, Tags map[string]string) error
}
