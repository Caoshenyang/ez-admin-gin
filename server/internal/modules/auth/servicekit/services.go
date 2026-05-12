package servicekit

import (
	authapp "ez-admin-gin/server/internal/modules/auth/application"
	authinfra "ez-admin-gin/server/internal/modules/auth/infra"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Services 聚合 auth 模块在 HTTP 层需要用到的应用服务，减少路由层重复拼装。
type Services struct {
	Login     *authapp.LoginService
	Refresh   *authapp.RefreshService
	Logout    *authapp.LogoutService
	Me        *authapp.MeService
	Account   *authapp.AccountService
	Menu      *authapp.MenuService
	Dashboard *authapp.DashboardService
}

type ServiceOptions struct {
	Config       *platformConfig.Config
	Log          *zap.Logger
	DB           *gorm.DB
	Redis        *goredis.Client
	Token        *authnPlatform.Manager
	RefreshStore *authnPlatform.RefreshTokenStore
}

// NewServices 收拢 auth 模块服务构造，统一 repository / transactor / token / redis 依赖装配。
func NewServices(opts ServiceOptions) Services {
	repo := authinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)

	return Services{
		Login:     authapp.NewLoginService(repo, opts.Token, nil, opts.Log),
		Refresh:   authapp.NewRefreshService(opts.RefreshStore, opts.Token),
		Logout:    authapp.NewLogoutService(opts.RefreshStore, opts.Token),
		Me:        authapp.NewMeService(),
		Account:   authapp.NewAccountService(transactor, repo),
		Menu:      authapp.NewMenuService(repo),
		Dashboard: authapp.NewDashboardService(opts.Config, opts.DB, repo, opts.Redis, opts.Log),
	}
}
