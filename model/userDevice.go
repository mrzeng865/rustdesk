package model

// UserDevice 用户-设备绑定表（用于限制账号可登录设备数量）
type UserDevice struct {
	IdModel

	UserId     uint   `json:"user_id" gorm:"not null;index;uniqueIndex:idx_user_device"`
	DeviceUuid string `json:"device_uuid" gorm:"default:'';not null;index;uniqueIndex:idx_user_device"` // RustDesk uuid
	DeviceId   string `json:"device_id" gorm:"default:'';not null;index"`                               // RustDesk ID

	DeviceName string `json:"device_name" gorm:"default:'';not null;"`
	DeviceOs   string `json:"device_os" gorm:"default:'';not null;"`
	DeviceType string `json:"device_type" gorm:"default:'';not null;"` // app/webclient/...

	FirstLoginAt int64  `json:"first_login_at" gorm:"default:0;not null;"`
	LastLoginAt  int64  `json:"last_login_at" gorm:"default:0;not null;"`
	LastLoginIp  string `json:"last_login_ip" gorm:"default:'';not null;"`

	TimeModel
}