package devices

type Repository interface {
	UpsertDevice(model DeviceModel) (*DeviceModel, error)
	GetDevice(uuid string, applicationID uint) (*DeviceModel, error)
	
	GetApplicationByUUID(uuid string) (*DeviceApplicationModel, error)
}
