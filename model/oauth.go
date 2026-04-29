package model

import (
	"errors"
	"strconv"
	"strings"
)

const OIDC_DEFAULT_SCOPES = "openid,profile,email"

const (
	// make sure the value shouldbe lowercase
	OauthTypeGithub  string = "github"
	OauthTypeGoogle  string = "google"
	OauthTypeOidc    string = "oidc"
	OauthTypeWebauth string = "webauth"
	OauthTypeLinuxdo string = "linuxdo"
	PKCEMethodS256   string = "S256"
	PKCEMethodPlain  string = "plain"
)

// Validate the oauth type
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
	//RedirectUrl  string `json:"redirect_url"`
	AutoRegister *bool  `json:"auto_register"`
	Scopes       string `json:"scopes"`
	Issuer       string `json:"issuer"`
	PkceEnable   *bool  `json:"pkce_enable"`
	PkceMethod   string `json:"pkce_method"`
	TimeModel
}

// Helper function to format oauth info, it's used in the update and create method
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
	// check if the op is empty, set the default value
	op := strings.TrimSpace(oa.Op)
	if op == "" && oauthType == OauthTypeOidc {
		oa.Op = OauthTypeOidc
	}
	// check the issuer, if the oauth type is google and the issuer is empty, set the issuer to the default value
	issuer := strings.TrimSpace(oa.Issuer)
	// If the oauth type is google and the issuer is empty, set the issuer to the default value
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

// OidcUser 接收 OIDC Provider 返回的 UserInfo JSON
type OidcUser struct {
	OauthUserBase
	Sub               string `json:"sub"`
	VerifiedEmail     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	Picture           string `json:"picture"`
	Nickname          string `json:"nickname"`          // 兼容飞书可能返回的 nickname 字段
	DisplayName       string `json:"display_name"`      // 兼容部分企业微信/钉钉
}

func (ou *OidcUser) ToOauthUser() *OauthUser {
	var username string
	
	// 优先级降级链
	if ou.PreferredUsername != "" {
		username = ou.PreferredUsername
	} else if ou.Name != "" {
		username = ou.Name
	} else if ou.Nickname != "" {
		username = ou.Nickname
	} else if ou.DisplayName != "" {
		username = ou.DisplayName
	} else if ou.Email != "" {
		username = strings.ToLower(ou.Email)
	} else {
		username = "oidc_" + ou.Sub
	}

	//  防重后缀：飞书多人同名会导致 uniqueIndex 冲突，自动加时间戳后缀保证唯一
	username = fmt.Sprintf("%s_%d", username, time.Now().UnixNano()%10000)

	//  临时调试日志：登录成功后会在 docker logs 中打印实际解析结果
	fmt.Printf("[OIDC DEBUG] sub=%s, name=%s, username_resolved=%s\n", ou.Sub, ou.Name, username)

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
		VerifiedEmail: true, // linux.do 用户邮箱默认已验证
		Picture:       lu.Avatar,
	}
}

type OauthList struct {
	Oauths []*Oauth `json:"list"`
	Pagination
}
