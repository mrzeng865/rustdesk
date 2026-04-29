package model

import (
	"errors"
	"strconv"
	"strings"
)

const OIDC_DEFAULT_SCOPES = "openid,profile,email"

const (
	OauthTypeGithub  string = "github"
	OauthTypeGoogle  string = "google"
	OauthTypeOidc    string = "oidc"
	OauthTypeWebauth string = "webauth"
	OauthTypeLinuxdo string = "linuxdo"
	PKCEMethodS256   string = "S256"
	PKCEMethodPlain  string = "plain"
)

func ValidateOauthType(oauthType string) error {
	switch oauthType {
	case OauthTypeGithub, OauthTypeGoogle, OauthTypeOidc, OauthTypeWebauth, OauthTypeLinuxdo:
		return nil
	default:
		return errors.New("invalid Oauth type")
	}
}

const (
	UserEndpointGithub  string = "https://api.github.com/user"
	UserEndpointLinuxdo string = "https://connect.linux.do/api/user"
	IssuerGoogle        string = "https://accounts.google.com"
)

type Oauth struct {
	IdModel
	Op           string `json:"op"`
	OauthType    string `json:"oauth_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AutoRegister *bool  `json:"auto_register"`
	Scopes       string `json:"scopes"`
	Issuer       string `json:"issuer"`
	PkceEnable   *bool  `json:"pkce_enable"`
	PkceMethod   string `json:"pkce_method"`
	TimeModel
}

func (oa *Oauth) FormatOauthInfo() error {
	oauthType := strings.TrimSpace(oa.OauthType)
	err := ValidateOauthType(oa.OauthType)
	if err != nil {
		return err
	}
	switch oauthType {
	case OauthTypeGithub:
		oa.Op = OauthTypeGithub
	case OauthTypeGoogle:
		oa.Op = OauthTypeGoogle
	case OauthTypeLinuxdo:
		oa.Op = OauthTypeLinuxdo
	}
	op := strings.TrimSpace(oa.Op)
	if op == "" && oauthType == OauthTypeOidc {
		oa.Op = OauthTypeOidc
	}
	issuer := strings.TrimSpace(oa.Issuer)
	if oauthType == OauthTypeGoogle && issuer == "" {
		oa.Issuer = IssuerGoogle
	}
	if oa.PkceEnable == nil {
		oa.PkceEnable = new(bool)
		*oa.PkceEnable = false
	}
	if oa.PkceMethod == "" {
		oa.PkceMethod = PKCEMethodS256
	}
	return nil
}

type OauthUser struct {
	OpenId        string `json:"open_id" gorm:"not null;index"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email,omitempty"`
	Picture       string `json:"picture,omitempty"`
}

func (ou *OauthUser) ToUser(user *User, overideUsername bool) {
	if overideUsername {
		user.Username = ou.Username
	}
	user.Email = ou.Email
	user.Nickname = ou.Name
	user.Avatar = ou.Picture
}

type OauthUserBase struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type OidcUser struct {
	OauthUserBase
	Sub               string `json:"sub"`
	VerifiedEmail     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	Picture           string `json:"picture"`
}

//  核心修改：飞书 OIDC 用户名降级映射逻辑
func (ou *OidcUser) ToOauthUser() *OauthUser {
	var username string
	// 优先级：标准字段 -> 飞书昵称 -> 邮箱 -> 兜底ID
	if ou.PreferredUsername != "" {
		username = ou.PreferredUsername
	} else if ou.Name != "" {
		username = ou.Name // 飞书默认返回昵称
	} else if ou.Email != "" {
		username = strings.ToLower(ou.Email)
	} else {
		username = "oidc_" + ou.Sub
	}

	return &OauthUser{
		OpenId:        ou.Sub,
		Name:          ou.Name,
		Username:      username,
		Email:         ou.Email,
		VerifiedEmail: ou.VerifiedEmail,
		Picture:       ou.Picture,
	}
}

type GithubUser struct {
	OauthUserBase
	Id            int    `json:"id"`
	Login         string `json:"login"`
	AvatarUrl     string `json:"avatar_url"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (gu *GithubUser) ToOauthUser() *OauthUser {
	username := strings.ToLower(gu.Login)
	return &OauthUser{
		OpenId:        strconv.Itoa(gu.Id),
		Name:          gu.Name,
		Username:      username,
		Email:         gu.Email,
		VerifiedEmail: gu.VerifiedEmail,
		Picture:       gu.AvatarUrl,
	}
}

type LinuxdoUser struct {
	OauthUserBase
	Id       int    `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar_url"`
}

func (lu *LinuxdoUser) ToOauthUser() *OauthUser {
	return &OauthUser{
		OpenId:        strconv.Itoa(lu.Id),
		Name:          lu.Name,
		Username:      strings.ToLower(lu.Username),
		Email:         lu.Email,
		VerifiedEmail: true,
		Picture:       lu.Avatar,
	}
}

type OauthList struct {
	Oauths []*Oauth `json:"list"`
	Pagination
}
