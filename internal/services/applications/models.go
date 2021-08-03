package applications

import (
	"time"
)

type ApplicationModel struct {
	ID                   uint
	UUID                 string
	Name                 string
	FCMSenderID          string
	FCMAdminJSON         string
	URL                  string
	AuthKey              string
	IdentityVerification bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	AccountID            uint `gorm:"index"`
}
