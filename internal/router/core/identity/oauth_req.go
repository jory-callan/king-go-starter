package identity

// OAuthAuthorizeReq OAuth 授权请求参数
type OAuthAuthorizeReq struct {
	ClientID     string `query:"client_id" validate:"required"`
	RedirectURI  string `query:"redirect_uri" validate:"required"`
	ResponseType string `query:"response_type" validate:"required,oneof=code token"`
	Scope        string `query:"scope"`
	State        string `query:"state"`
}

// OAuthTokenReq OAuth 获取令牌请求参数
type OAuthTokenReq struct {
	ClientID     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	GrantType    string `json:"grant_type" validate:"required,oneof=authorization_code refresh_token password client_credentials"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// OAuthUserInfoReq OAuth 获取用户信息请求参数
type OAuthUserInfoReq struct {
	AccessToken string `header:"Authorization" validate:"required"`
}
