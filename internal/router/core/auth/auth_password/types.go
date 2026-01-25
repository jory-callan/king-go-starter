package auth_password

const (
	AuthTypePassword string = "password" // 密码认证
	AuthTypeEmail    string = "email"    // 邮箱认证
	AuthTypePhone    string = "phone"    // 手机认证
	AuthType2FA      string = "2fa"      // 两步验证
	AuthTypeOAuth2   string = "oauth2"   // OAuth2 认证
	AuthTypeSSO      string = "sso"      // 单点登录
)

const (
	LoginTypeSuccess string = "success" // 登录成功
	LoginTypeFailed  string = "failed"  // 登录失败
)
