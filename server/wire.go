//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/mazrean/Quantainer/auth"
	traq "github.com/mazrean/Quantainer/auth/traQ"
	bot "github.com/mazrean/Quantainer/bot/traq"
	"github.com/mazrean/Quantainer/cache"
	"github.com/mazrean/Quantainer/cache/ristretto"
	v1Handler "github.com/mazrean/Quantainer/handler/v1"
	"github.com/mazrean/Quantainer/pkg/common"
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/repository/gorm2"
	"github.com/mazrean/Quantainer/service"
	v1Service "github.com/mazrean/Quantainer/service/v1"
	"github.com/mazrean/Quantainer/storage"
	"github.com/mazrean/Quantainer/storage/local"
	"github.com/mazrean/Quantainer/storage/swift"
)

type Config struct {
	IsProduction      common.IsProduction
	SessionKey        common.SessionKey
	SessionSecret     common.SessionSecret
	TraQBaseURL       common.TraQBaseURL
	OAuthClientID     common.ClientID
	SwiftAuthURL      common.SwiftAuthURL
	SwiftUserName     common.SwiftUserName
	SwiftPassword     common.SwiftPassword
	SwiftTenantID     common.SwiftTenantID
	SwiftTenantName   common.SwiftTenantName
	SwiftContainer    common.SwiftContainer
	FilePath          common.FilePath
	AccessToken       common.AccessToken
	VerificationToken common.VerificationToken
	DefaultChannels   common.DefaultChannels
	UpdatedAt         common.UpdatedAt
	HttpClient        *http.Client
}

type Storage struct {
	File storage.File
}

func newStorage(file storage.File) *Storage {
	return &Storage{
		File: file,
	}
}

var (
	isProductionField      = wire.FieldsOf(new(*Config), "IsProduction")
	sessionKeyField        = wire.FieldsOf(new(*Config), "SessionKey")
	sessionSecretField     = wire.FieldsOf(new(*Config), "SessionSecret")
	traQBaseURLField       = wire.FieldsOf(new(*Config), "TraQBaseURL")
	oAuthClientIDField     = wire.FieldsOf(new(*Config), "OAuthClientID")
	swiftAuthURLField      = wire.FieldsOf(new(*Config), "SwiftAuthURL")
	swiftUserNameField     = wire.FieldsOf(new(*Config), "SwiftUserName")
	swiftPasswordField     = wire.FieldsOf(new(*Config), "SwiftPassword")
	swiftTenantIDField     = wire.FieldsOf(new(*Config), "SwiftTenantID")
	swiftTenantNameField   = wire.FieldsOf(new(*Config), "SwiftTenantName")
	swiftContainerField    = wire.FieldsOf(new(*Config), "SwiftContainer")
	filePathField          = wire.FieldsOf(new(*Config), "FilePath")
	accessTokenField       = wire.FieldsOf(new(*Config), "AccessToken")
	verificationTokenField = wire.FieldsOf(new(*Config), "VerificationToken")
	defaultChannelsField   = wire.FieldsOf(new(*Config), "DefaultChannels")
	updatedAtField         = wire.FieldsOf(new(*Config), "UpdatedAt")
	httpClientField        = wire.FieldsOf(new(*Config), "HttpClient")
)

func injectedStorage(config *Config) (*Storage, error) {
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
		wire.Bind(new(storage.File), new(*swift.File)),
		swift.NewClient,
		swift.NewFile,
		newStorage,
	)

	return nil, nil
}

func injectLocalStorage(config *Config) (*Storage, error) {
	wire.Build(
		filePathField,
		wire.Bind(new(storage.File), new(*local.File)),
		local.NewDirectoryManager,
		local.NewFile,
		newStorage,
	)

	return nil, nil
}

var (
	dbBind                      = wire.Bind(new(repository.DB), new(*gorm2.DB))
	fileRepositoryBind          = wire.Bind(new(repository.File), new(*gorm2.File))
	resourceRepositoryBind      = wire.Bind(new(repository.Resource), new(*gorm2.Resource))
	groupRepositoryBind         = wire.Bind(new(repository.Group), new(*gorm2.Group))
	administratorRepositoryBind = wire.Bind(new(repository.Administrator), new(*gorm2.Administrator))

	oidcAuthBind = wire.Bind(new(auth.OIDC), new(*traq.OIDC))
	userAuthBind = wire.Bind(new(auth.User), new(*traq.User))

	userCacheBind = wire.Bind(new(cache.User), new(*ristretto.User))

	oidcServiceBind     = wire.Bind(new(service.OIDC), new(*v1Service.OIDC))
	userServiceBind     = wire.Bind(new(service.User), new(*v1Service.User))
	fileServiceBind     = wire.Bind(new(service.File), new(*v1Service.File))
	resourceServiceBind = wire.Bind(new(service.Resource), new(*v1Service.Resource))
	groupServiceBind    = wire.Bind(new(service.Group), new(*v1Service.Group))

	fileField = wire.FieldsOf(new(*Storage), "File")
)

type Service struct {
	*v1Handler.API
	*bot.Bot
}

func NewService(api *v1Handler.API, b *bot.Bot) *Service {
	return &Service{
		API: api,
		Bot: b,
	}
}

func InjectService(config *Config) (*Service, error) {
	wire.Build(
		isProductionField,
		sessionKeyField,
		sessionSecretField,
		traQBaseURLField,
		oAuthClientIDField,
		httpClientField,
		fileField,
		accessTokenField,
		verificationTokenField,
		defaultChannelsField,
		updatedAtField,
		dbBind,
		fileRepositoryBind,
		resourceRepositoryBind,
		groupRepositoryBind,
		administratorRepositoryBind,
		oidcAuthBind,
		userAuthBind,
		userCacheBind,
		oidcServiceBind,
		userServiceBind,
		fileServiceBind,
		resourceServiceBind,
		groupServiceBind,
		gorm2.NewDB,
		gorm2.NewFile,
		gorm2.NewResource,
		gorm2.NewGroup,
		gorm2.NewAdministrator,
		traq.NewOIDC,
		traq.NewUser,
		ristretto.NewUser,
		v1Service.NewOIDC,
		v1Service.NewUser,
		v1Service.NewUserUtils,
		v1Service.NewFile,
		v1Service.NewResource,
		v1Service.NewGroup,
		v1Handler.NewAPI,
		v1Handler.NewSession,
		v1Handler.NewOAuth2,
		v1Handler.NewUser,
		v1Handler.NewChecker,
		v1Handler.NewFile,
		v1Handler.NewResource,
		v1Handler.NewGroup,
		bot.NewBot,
		injectedStorage,
		NewService,
	)
	return nil, nil
}
