package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	req "github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

type UserDevice struct{}

type UserDeviceItem struct {
	DeviceUuid    string `json:"device_uuid"`
	DeviceId      string `json:"device_id"`
	DeviceName    string `json:"device_name"`
	DeviceOs      string `json:"device_os"`
	DeviceType    string `json:"device_type"`
	FirstLoginAt  int64  `json:"first_login_at"`
	LastLoginAt   int64  `json:"last_login_at"`
	LastLoginIp   string `json:"last_login_ip"`

	PeerHostname       string `json:"peer_hostname"`
	PeerVersion        string `json:"peer_version"`
	PeerLastOnlineTime int64  `json:"peer_last_online_time"`
	PeerLastOnlineIp   string `json:"peer_last_online_ip"`
}

type UserDeviceListRes struct {
	UserId     uint             `json:"user_id"`
	MaxDevices uint             `json:"max_devices"`
	List       []UserDeviceItem `json:"list"`
}

// List 用户设备绑定列表（带 peer 信息）
func (ct *UserDevice) List(c *gin.Context) {
	userIdStr := c.Query("user_id")
	uid, _ := strconv.Atoi(userIdStr)
	if uid <= 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}

	u := service.AllService.UserService.InfoById(uint(uid))
	if u == nil || u.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}

	var uds []model.UserDevice
	service.DB.Where("user_id = ?", u.Id).Order("last_login_at desc").Find(&uds)

	// 批量查 peer
	uuids := make([]string, 0, len(uds))
	for _, d := range uds {
		if d.DeviceUuid != "" {
			uuids = append(uuids, d.DeviceUuid)
		}
	}
	peerMap := map[string]*model.Peer{}
	if len(uuids) > 0 {
		var peers []model.Peer
		service.DB.Where("uuid in (?)", uuids).Find(&peers)
		for i := range peers {
			p := peers[i]
			peerMap[p.Uuid] = &p
		}
	}

	items := make([]UserDeviceItem, 0, len(uds))
	for _, d := range uds {
		item := UserDeviceItem{
			DeviceUuid:   d.DeviceUuid,
			DeviceId:     d.DeviceId,
			DeviceName:   d.DeviceName,
			DeviceOs:     d.DeviceOs,
			DeviceType:   d.DeviceType,
			FirstLoginAt: d.FirstLoginAt,
			LastLoginAt:  d.LastLoginAt,
			LastLoginIp:  d.LastLoginIp,
		}
		if p, ok := peerMap[d.DeviceUuid]; ok && p != nil {
			item.PeerHostname = p.Hostname
			item.PeerVersion = p.Version
			item.PeerLastOnlineTime = p.LastOnlineTime
			item.PeerLastOnlineIp = p.LastOnlineIp
		}
		items = append(items, item)
	}

	response.Success(c, UserDeviceListRes{
		UserId:     u.Id,
		MaxDevices: func() uint { if u.MaxDevices == 0 { return 1 }; return u.MaxDevices }(),
		List:       items,
	})
}

// SetLimit 设置最大设备数
func (ct *UserDevice) SetLimit(c *gin.Context) {
	f := &req.UserDeviceSetLimitForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}

	if err := service.AllService.UserDeviceService.SetMaxDevices(f.UserId, f.MaxDevices); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Unbind 解绑设备
func (ct *UserDevice) Unbind(c *gin.Context) {
	f := &req.UserDeviceUnbindForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}

	if err := service.AllService.UserDeviceService.UnbindDevice(f.UserId, f.DeviceUuid); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}