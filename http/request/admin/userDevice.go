package admin

type UserDeviceQuery struct {
	UserId uint `form:"user_id" json:"user_id" validate:"required"`
}

type UserDeviceSetLimitForm struct {
	UserId     uint `json:"user_id" validate:"required"`
	MaxDevices uint `json:"max_devices" validate:"required,gte=1"`
}

type UserDeviceUnbindForm struct {
	UserId     uint   `json:"user_id" validate:"required"`
	DeviceUuid string `json:"device_uuid" validate:"required"`
}