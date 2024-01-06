package backgroundsvcs

import (
	"github.com/xquare-dashboard/pkg/api"
	"github.com/xquare-dashboard/pkg/infra/metrics"
	"github.com/xquare-dashboard/pkg/infra/remotecache"
	"github.com/xquare-dashboard/pkg/infra/tracing"
	uss "github.com/xquare-dashboard/pkg/infra/usagestats/service"
	"github.com/xquare-dashboard/pkg/infra/usagestats/statscollector"
	"github.com/xquare-dashboard/pkg/registry"
	apiregistry "github.com/xquare-dashboard/pkg/registry/apis"
	"github.com/xquare-dashboard/pkg/services/alerting"
	"github.com/xquare-dashboard/pkg/services/anonymous/anonimpl"
	"github.com/xquare-dashboard/pkg/services/auth"
	"github.com/xquare-dashboard/pkg/services/cleanup"
	"github.com/xquare-dashboard/pkg/services/dashboardsnapshots"
	grafanaapiserver "github.com/xquare-dashboard/pkg/services/grafana-apiserver"
	"github.com/xquare-dashboard/pkg/services/grpcserver"
	"github.com/xquare-dashboard/pkg/services/guardian"
	ldapapi "github.com/xquare-dashboard/pkg/services/ldap/api"
	"github.com/xquare-dashboard/pkg/services/live"
	"github.com/xquare-dashboard/pkg/services/live/pushhttp"
	"github.com/xquare-dashboard/pkg/services/loginattempt/loginattemptimpl"
	"github.com/xquare-dashboard/pkg/services/ngalert"
	"github.com/xquare-dashboard/pkg/services/notifications"
	plugindashboardsservice "github.com/xquare-dashboard/pkg/services/plugindashboards/service"
	"github.com/xquare-dashboard/pkg/services/pluginsintegration/angulardetectorsprovider"
	"github.com/xquare-dashboard/pkg/services/pluginsintegration/keyretriever/dynamic"
	pluginStore "github.com/xquare-dashboard/pkg/services/pluginsintegration/pluginstore"
	"github.com/xquare-dashboard/pkg/services/provisioning"
	publicdashboardsmetric "github.com/xquare-dashboard/pkg/services/publicdashboards/metric"
	"github.com/xquare-dashboard/pkg/services/rendering"
	"github.com/xquare-dashboard/pkg/services/searchV2"
	secretsMigrations "github.com/xquare-dashboard/pkg/services/secrets/kvstore/migrations"
	secretsManager "github.com/xquare-dashboard/pkg/services/secrets/manager"
	"github.com/xquare-dashboard/pkg/services/serviceaccounts"
	samanager "github.com/xquare-dashboard/pkg/services/serviceaccounts/manager"
	"github.com/xquare-dashboard/pkg/services/ssosettings"
	"github.com/xquare-dashboard/pkg/services/store"
	"github.com/xquare-dashboard/pkg/services/store/entity"
	"github.com/xquare-dashboard/pkg/services/store/sanitizer"
	"github.com/xquare-dashboard/pkg/services/supportbundles/supportbundlesimpl"
	"github.com/xquare-dashboard/pkg/services/team/teamapi"
	"github.com/xquare-dashboard/pkg/services/updatechecker"
)

func ProvideBackgroundServiceRegistry(
	httpServer *api.HTTPServer, ng *ngalert.AlertNG, cleanup *cleanup.CleanUpService, live *live.GrafanaLive,
	pushGateway *pushhttp.Gateway, notifications *notifications.NotificationService, pluginStore *pluginStore.Service,
	rendering *rendering.RenderingService, tokenService auth.UserTokenBackgroundService, tracing *tracing.TracingService,
	provisioning *provisioning.ProvisioningServiceImpl, alerting *alerting.AlertEngine, usageStats *uss.UsageStats,
	statsCollector *statscollector.Service, grafanaUpdateChecker *updatechecker.GrafanaService,
	pluginsUpdateChecker *updatechecker.PluginsService, metrics *metrics.InternalMetricsService,
	secretsService *secretsManager.SecretsService, remoteCache *remotecache.RemoteCache, StorageService store.StorageService, searchService searchV2.SearchService, entityEventsService store.EntityEventsService,
	saService *samanager.ServiceAccountsService, grpcServerProvider grpcserver.Provider,
	secretMigrationProvider secretsMigrations.SecretMigrationProvider, loginAttemptService *loginattemptimpl.Service,
	bundleService *supportbundlesimpl.Service, publicDashboardsMetric *publicdashboardsmetric.Service,
	keyRetriever *dynamic.KeyRetriever, dynamicAngularDetectorsProvider *angulardetectorsprovider.Dynamic,
	grafanaAPIServer grafanaapiserver.Service,
	anon *anonimpl.AnonDeviceService,
	// Need to make sure these are initialized, is there a better place to put them?
	_ dashboardsnapshots.Service, _ *alerting.AlertNotificationService,
	_ serviceaccounts.Service, _ *guardian.Provider,
	_ *plugindashboardsservice.DashboardUpdater, _ *sanitizer.Provider,
	_ *grpcserver.HealthService, _ entity.EntityStoreServer, _ *grpcserver.ReflectionService, _ *ldapapi.Service,
	_ *apiregistry.Service, _ auth.IDService, _ *teamapi.TeamAPI, _ ssosettings.Service,
) *BackgroundServiceRegistry {
	return NewBackgroundServiceRegistry(
		httpServer,
		ng,
		cleanup,
		live,
		pushGateway,
		notifications,
		rendering,
		tokenService,
		provisioning,
		alerting,
		grafanaUpdateChecker,
		pluginsUpdateChecker,
		metrics,
		usageStats,
		statsCollector,
		tracing,
		remoteCache,
		secretsService,
		StorageService,
		searchService,
		entityEventsService,
		grpcServerProvider,
		saService,
		pluginStore,
		secretMigrationProvider,
		loginAttemptService,
		bundleService,
		publicDashboardsMetric,
		keyRetriever,
		dynamicAngularDetectorsProvider,
		grafanaAPIServer,
		anon,
	)
}

// BackgroundServiceRegistry provides background services.
type BackgroundServiceRegistry struct {
	Services []registry.BackgroundService
}

func NewBackgroundServiceRegistry(services ...registry.BackgroundService) *BackgroundServiceRegistry {
	return &BackgroundServiceRegistry{services}
}

func (r *BackgroundServiceRegistry) GetServices() []registry.BackgroundService {
	return r.Services
}
