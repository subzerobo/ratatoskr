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

// Cache Data Models

type ApplicationCachedDataModel struct {
	ChannelList              []ChannelListDataModel `json:"chnl_lst"`
	UseIdentityVerification  bool                   `json:"use_email_auth"`
	FirebaseAnalytics        bool                   `json:"fba"`
	FCMId                    string                 `json:"android_sender_id"`
	CleanGroupOnSummaryClick bool                   `json:"clear_group_on_summary_click"`
	ReceiveReceipt           bool                   `json:"receive_receipts_enable"`
}

type ChannelListDataModel struct {
	Channel          ChannelDataModel `json:"chnl"`
	Sound            string           `json:"sound,omitempty"`
	VibrationPattern string           `json:"vib_pt,omitempty"`
	Priority         string           `json:"pri,omitempty"`
	Badge            string           `json:"bdg,omitempty"`
	LedColor         string           `json:"led_color,omitempty"`
	LockScreen       int              `json:"vis,omitempty"`
}

type ChannelDataModel struct {
	ID          string `json:"id"`
	Name        string `json:"nm"`
	Description string `json:"dscr"`
	GroupID     string `json:"grp_id"`
	GroupName   string `json:"grp_nm"`
}
