package service

import (
	"errors"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
)

var (
	// i18n message id：在 resources/i18n/*.toml 里加同名 key
	ErrDeviceLimitExceeded = errors.New("DeviceLimitExceeded")
	ErrParamsError         = errors.New("ParamsError")
)

type UserDeviceService struct{}

// EnsureBoundOrReject：登录前检查设备是否允许登录；必要时写入绑定记录
func (s *UserDeviceService) EnsureBoundOrReject(
	u *model.User,
	uuid string,
	deviceId string,
	deviceName string,
	deviceOs string,
	deviceType string,
	ip string,
) error {
	if u == nil || u.Id == 0 {
		return ErrParamsError
	}
	if uuid == "" {
		// 没有 uuid 无法做绑定限制，直接按参数错误处理（更安全）
		return ErrParamsError
	}

	max := u.MaxDevices
	if max == 0 {
		max = 1
	}

	now := time.Now().Unix()

	// 已绑定？更新 last_login 并放行
	exist := &model.UserDevice{}
	DB.Where("user_id = ? and device_uuid = ?", u.Id, uuid).First(exist)
	if exist.Id != 0 {
		DB.Model(exist).Updates(map[string]interface{}{
			"device_id":      deviceId,
			"device_name":    deviceName,
			"device_os":      deviceOs,
			"device_type":    deviceType,
			"last_login_at":  now,
			"last_login_ip":  ip,
		})
		return nil
	}

	// 未绑定：检查当前绑定数量
	var cnt int64
	DB.Model(&model.UserDevice{}).Where("user_id = ?", u.Id).Count(&cnt)
	if uint(cnt) >= max {
		return ErrDeviceLimitExceeded
	}

	// 允许绑定：创建绑定记录
	ud := &model.UserDevice{
		UserId:       u.Id,
		DeviceUuid:   uuid,
		DeviceId:     deviceId,
		DeviceName:   deviceName,
		DeviceOs:     deviceOs,
		DeviceType:   deviceType,
		FirstLoginAt: now,
		LastLoginAt:  now,
		LastLoginIp:  ip,
	}
	return DB.Create(ud).Error
}

// SetMaxDevices：设置用户最大设备数（最小 1）
func (s *UserDeviceService) SetMaxDevices(userId uint, max uint) error {
	if userId == 0 || max == 0 {
		return ErrParamsError
	}
	return DB.Model(&model.User{}).Where("id = ?", userId).Update("max_devices", max).Error
}

// UnbindDevice：解绑设备 + 清除该 uuid 的 token + 解绑 peer.user_id
func (s *UserDeviceService) UnbindDevice(userId uint, uuid string) error {
	if userId == 0 || uuid == "" {
		return ErrParamsError
	}

	// 删除绑定
	if err := DB.Where("user_id = ? and device_uuid = ?", userId, uuid).Delete(&model.UserDevice{}).Error; err != nil {
		return err
	}

	// 清除 token（注意：如果你启用的是 JWT 且服务端只做签名校验，则已发出的 JWT 仍可能在过期前有效）
	_ = AllService.UserService.FlushTokenByUuid(uuid)

	// 解绑 peer.user_id（可选，但建议做，后台“设备管理”列表也会更干净）
	AllService.PeerService.UuidUnbindUserId(uuid, userId)

	return nil
}