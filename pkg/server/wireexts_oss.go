//go:build wireinject && oss
// +build wireinject,oss

// This file should contain wiresets which contain OSS-specific implementations.
package server

import (
	"github.com/google/wire"

	"github.com/xquare-dashboard/pkg/infra/metrics"
	"github.com/xquare-dashboard/pkg/plugins"
	"github.com/xquare-dashboard/pkg/plugins/manager"
	"github.com/xquare-dashboard/pkg/registry"
	"github.com/xquare-dashboard/pkg/registry/backgroundsvcs"
	"github.com/xquare-dashboard/pkg/registry/usagestatssvcs"
	"github.com/xquare-dashboard/pkg/services/accesscontrol"
	"github.com/xquare-dashboard/pkg/services/accesscontrol/acimpl"
	"github.com/xquare-dashboard/pkg/services/accesscontrol/ossaccesscontrol"
	"github.com/xquare-dashboard/pkg/services/anonymous"
	"github.com/xquare-dashboard/pkg/services/anonymous/anonimpl"
	"github.com/xquare-dashboard/pkg/services/auth"
	"github.com/xquare-dashboard/pkg/services/auth/authimpl"
	"github.com/xquare-dashboard/pkg/services/auth/idimpl"
	"github.com/xquare-dashboard/pkg/services/caching"
	"github.com/xquare-dashboard/pkg/services/datasources/guardian"
	"github.com/xquare-dashboard/pkg/services/encryption"
	encryptionprovider "github.com/xquare-dashboard/pkg/services/encryption/provider"
	"github.com/xquare-dashboard/pkg/services/featuremgmt"
	"github.com/xquare-dashboard/pkg/services/hooks"
	"github.com/xquare-dashboard/pkg/services/kmsproviders"
	"github.com/xquare-dashboard/pkg/services/kmsproviders/osskmsproviders"
	"github.com/xquare-dashboard/pkg/services/ldap"
	"github.com/xquare-dashboard/pkg/services/licensing"
	"github.com/xquare-dashboard/pkg/services/login"
	"github.com/xquare-dashboard/pkg/services/login/authinfoimpl"
	"github.com/xquare-dashboard/pkg/services/pluginsintegration"
	"github.com/xquare-dashboard/pkg/services/provisioning"
	"github.com/xquare-dashboard/pkg/services/publicdashboards"
	publicdashboardsApi "github.com/xquare-dashboard/pkg/services/publicdashboards/api"
	publicdashboardsService "github.com/xquare-dashboard/pkg/services/publicdashboards/service"
	"github.com/xquare-dashboard/pkg/services/searchusers"
	"github.com/xquare-dashboard/pkg/services/searchusers/filters"
	"github.com/xquare-dashboard/pkg/services/secrets"
	secretsMigrator "github.com/xquare-dashboard/pkg/services/secrets/migrator"
	"github.com/xquare-dashboard/pkg/services/sqlstore/migrations"
	"github.com/xquare-dashboard/pkg/services/user"
	"github.com/xquare-dashboard/pkg/services/validations"
	"github.com/xquare-dashboard/pkg/setting"
)

var wireExtsBasicSet = wire.NewSet(
	authimpl.ProvideUserAuthTokenService,
	wire.Bind(new(auth.UserTokenService), new(*authimpl.UserAuthTokenService)),
	wire.Bind(new(auth.UserTokenBackgroundService), new(*authimpl.UserAuthTokenService)),
	anonimpl.ProvideAnonymousDeviceService,
	wire.Bind(new(anonymous.Service), new(*anonimpl.AnonDeviceService)),
	licensing.ProvideService,
	wire.Bind(new(licensing.Licensing), new(*licensing.OSSLicensingService)),
	setting.ProvideProvider,
	wire.Bind(new(setting.Provider), new(*setting.OSSImpl)),
	acimpl.ProvideService,
	wire.Bind(new(accesscontrol.RoleRegistry), new(*acimpl.Service)),
	wire.Bind(new(plugins.RoleRegistry), new(*acimpl.Service)),
	wire.Bind(new(accesscontrol.Service), new(*acimpl.Service)),
	validations.ProvideValidator,
	wire.Bind(new(validations.PluginRequestValidator), new(*validations.OSSPluginRequestValidator)),
	provisioning.ProvideService,
	wire.Bind(new(provisioning.ProvisioningService), new(*provisioning.ProvisioningServiceImpl)),
	backgroundsvcs.ProvideBackgroundServiceRegistry,
	wire.Bind(new(registry.BackgroundServiceRegistry), new(*backgroundsvcs.BackgroundServiceRegistry)),
	migrations.ProvideOSSMigrations,
	wire.Bind(new(registry.DatabaseMigrator), new(*migrations.OSSMigrations)),
	authinfoimpl.ProvideOSSUserProtectionService,
	wire.Bind(new(login.UserProtectionService), new(*authinfoimpl.OSSUserProtectionImpl)),
	encryptionprovider.ProvideEncryptionProvider,
	wire.Bind(new(encryption.Provider), new(encryptionprovider.Provider)),
	filters.ProvideOSSSearchUserFilter,
	wire.Bind(new(user.SearchUserFilter), new(*filters.OSSSearchUserFilter)),
	searchusers.ProvideUsersService,
	wire.Bind(new(searchusers.Service), new(*searchusers.OSSService)),
	osskmsproviders.ProvideService,
	wire.Bind(new(kmsproviders.Service), new(osskmsproviders.Service)),
	ldap.ProvideGroupsService,
	wire.Bind(new(ldap.Groups), new(*ldap.OSSGroups)),
	guardian.ProvideGuardian,
	wire.Bind(new(guardian.DatasourceGuardianProvider), new(*guardian.OSSProvider)),
	usagestatssvcs.ProvideUsageStatsProvidersRegistry,
	wire.Bind(new(registry.UsageStatsProvidersRegistry), new(*usagestatssvcs.UsageStatsProvidersRegistry)),
	ossaccesscontrol.ProvideDatasourcePermissionsService,
	wire.Bind(new(accesscontrol.DatasourcePermissionsService), new(*ossaccesscontrol.DatasourcePermissionsService)),
	pluginsintegration.WireExtensionSet,
	publicdashboardsApi.ProvideMiddleware,
	wire.Bind(new(publicdashboards.Middleware), new(*publicdashboardsApi.Middleware)),
	publicdashboardsService.ProvideServiceWrapper,
	wire.Bind(new(publicdashboards.ServiceWrapper), new(*publicdashboardsService.PublicDashboardServiceWrapperImpl)),
	caching.ProvideCachingService,
	wire.Bind(new(caching.CachingService), new(*caching.OSSCachingService)),
	secretsMigrator.ProvideSecretsMigrator,
	wire.Bind(new(secrets.Migrator), new(*secretsMigrator.SecretsMigrator)),
	idimpl.ProvideLocalSigner,
	wire.Bind(new(auth.IDSigner), new(*idimpl.LocalSigner)),
	manager.ProvideInstaller,
	wire.Bind(new(plugins.Installer), new(*manager.PluginInstaller)),
)

var wireExtsSet = wire.NewSet(
	wireSet,
	wireExtsBasicSet,
)

var wireExtsCLISet = wire.NewSet(
	wireCLISet,
	wireExtsBasicSet,
)

var wireExtsTestSet = wire.NewSet(
	wireTestSet,
	wireExtsBasicSet,
)

// The wireExtsBaseCLISet is a simplified set of dependencies for the OSS CLI,
// suitable for running background services and targeted dskit modules without
// starting up the full Grafana server.
var wireExtsBaseCLISet = wire.NewSet(
	NewModuleRunner,

	metrics.WireSet,
	featuremgmt.ProvideManagerService,
	featuremgmt.ProvideToggles,
	hooks.ProvideService,
	setting.ProvideProvider, wire.Bind(new(setting.Provider), new(*setting.OSSImpl)),
	licensing.ProvideService, wire.Bind(new(licensing.Licensing), new(*licensing.OSSLicensingService)),
)

// wireModuleServerSet is a wire set for the ModuleServer.
var wireExtsModuleServerSet = wire.NewSet(
	NewModule,
	wireExtsBaseCLISet,
)
