package applications

import (
	"time"
)

type ApplicationModel struct {
	ID                   uint
	AccountID            uint
	UUID                 string
	Name                 string
	FCMSenderID          string
	FCMAdminJSON         string
	URL                  string
	AuthKey              string
	IdentityVerification bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type AndroidGroupModel struct {
	ID            uint                        `json:"id"`
	ApplicationID uint                        `json:"-"`
	GroupName     string                      `json:"group_name"`
	GroupUUID     string                      `json:"group_uuid"`
	CreatedAt     time.Time                   `json:"-"`
	UpdatedAt     time.Time                   `json:"-"`
	Categories    []AndroidGroupCategoryModel `json:"categories"`
}

type AndroidGroupCategoryModel struct {
	ID                  uint      `json:"id"`
	CategoryUUID        string    `json:"uuid"`
	CategoryName        string    `json:"name"`
	CategoryDescription string    `json:"description"`
	Priority            string    `json:"priority"`
	Sound               int       `json:"sound"`
	SoundName           string    `json:"sound_name"`
	Vibration           int       `json:"vibration"`
	VibrationPattern    string    `json:"vibration_pattern "`
	Led                 int       `json:"led"`
	LedColor            string    `json:"led_color"`
	EnableBadge         int       `json:"enable_badge"`
	LockScreen          int       `json:"lock_screen"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
	AndroidGroupID      uint      `json:"-"`
}
