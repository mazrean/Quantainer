//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/mazrean/Quantainer/auth"
	traq "github.com/mazrean/Quantainer/auth/traQ"
	"github.com/mazrean/Quantainer/cache"
	"github.com/mazrean/Quantainer/cache/ristretto"
	v1Handler "github.com/mazrean/Quantainer/handler/v1"
	"github.com/mazrean/Quantainer/pkg/common"
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/repository/gorm2"
	"github.com/mazrean/Quantainer/service"
	v1Service "github.com/mazrean/Quantainer/service/v1"
)

type Config struct {
	IsProduction    common.IsProduction
	SessionKey      common.SessionKey
	SessionSecret   common.SessionSecret
	TraQBaseURL     common.TraQBaseURL
	OAuthClientID   common.ClientID
	SwiftAuthURL    common.SwiftAuthURL
	SwiftUserName   common.SwiftUserName
	SwiftPassword   common.SwiftPassword
	SwiftTenantID   common.SwiftTenantID
	SwiftTenantName common.SwiftTenantName
	SwiftContainer  common.SwiftContainer
	FilePath        common.FilePath
	HttpClient      *http.Client
}

type Storage struct {
}

func newStorage() *Storage {
	return &Storage{}
}

var (
	isProductionField    = wire.FieldsOf(new(*Config), "IsProduction")
	sessionKeyField      = wire.FieldsOf(new(*Config), "SessionKey")
	sessionSecretField   = wire.FieldsOf(new(*Config), "SessionSecret")
	traQBaseURLField     = wire.FieldsOf(new(*Config), "TraQBaseURL")
	oAuthClientIDField   = wire.FieldsOf(new(*Config), "OAuthClientID")
	swiftAuthURLField    = wire.FieldsOf(new(*Config), "SwiftAuthURL")
	swiftUserNameField   = wire.FieldsOf(new(*Config), "SwiftUserName")
	swiftPasswordField   = wire.FieldsOf(new(*Config), "SwiftPassword")
	swiftTenantIDField   = wire.FieldsOf(new(*Config), "SwiftTenantID")
	swiftTenantNameField = wire.FieldsOf(new(*Config), "SwiftTenantName")
	swiftContainerField  = wire.FieldsOf(new(*Config), "SwiftContainer")
	filePathField        = wire.FieldsOf(new(*Config), "FilePath")
	httpClientField      = wire.FieldsOf(new(*Config), "HttpClient")
)

/*func injectedStorage(config *Config) (*Storage, error) {
	if config.IsProduction {
		return injectSwiftStorage(config)
	}

	return injectLocalStorage(config)
}

func injectSwiftStorage(config *Config) (*Storage, error) {
	wire.Build(
		swiftAuthURLField,
		swiftUserNameField,
		swiftPasswordField,
		swiftTenantIDField,
		swiftTenantNameField,
		swiftContainerField,
		filePathField,
		swift.NewClient,
		newStorage,
	)

	return nil, nil
}

func injectLocalStorage(config *Config) (*Storage, error) {
	wire.Build(
		filePathField,
		local.NewDirectoryManager,
		newStorage,
	)

	return nil, nil
}*/

var (
	dbBind = wire.Bind(new(repository.DB), new(*gorm2.DB))

	oidcAuthBind = wire.Bind(new(auth.OIDC), new(*traq.OIDC))
	userAuthBind = wire.Bind(new(auth.User), new(*traq.User))

	userCacheBind = wire.Bind(new(cache.User), new(*ristretto.User))

	oidcServiceBind = wire.Bind(new(service.OIDC), new(*v1Service.OIDC))
	userServiceBind = wire.Bind(new(service.User), new(*v1Service.User))
)

func InjectAPI(config *Config) (*v1Handler.API, error) {
	wire.Build(
		//isProductionField,
		sessionKeyField,
		sessionSecretField,
		traQBaseURLField,
		oAuthClientIDField,
		httpClientField,
		// dbBind,
		oidcAuthBind,
		userAuthBind,
		userCacheBind,
		oidcServiceBind,
		userServiceBind,
		//gorm2.NewDB,
		traq.NewOIDC,
		traq.NewUser,
		ristretto.NewUser,
		v1Service.NewOIDC,
		v1Service.NewUser,
		v1Service.NewUserUtils,
		v1Handler.NewAPI,
		v1Handler.NewSession,
		v1Handler.NewOAuth2,
		v1Handler.NewUser,
		v1Handler.NewChecker,
		// injectedStorage,
	)
	return nil, nil
}
