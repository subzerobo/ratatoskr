package devices

import "time"

type DeviceModel struct {
	ID                 uint
	UUID               string
	DeviceType         string
	Identifier         string
	Language           string
	Timezone           int
	AppVersion         string
	DeviceVendor       string
	DeviceModel        string
	DeviceOS           string
	DeviceOSVersion    string
	ADID               string
	SDK                string
	SessionCount       int
	NotificationTypes  int
	Long               float32
	Lat                float32
	Country            string
	ExternalUserID     string
	ExternalUserIDHash string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ApplicationID      uint
	Tags               map[string]string
}

type DeviceApplicationModel struct {
	ID                   uint
	UUID                 string
	AuthKey              string
	IdentityVerification bool
	AccountID            uint `gorm:"index"`
}
