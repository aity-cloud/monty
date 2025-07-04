package extensions

/*
Contains the AlertManager Monty embedded server implementation.
The embedded service must be run within the same process as each
deploymed node in the AlertManager cluster.
*/

import (
	"context"
	"errors"
	"net/http"

	"github.com/aity-cloud/monty/pkg/alerting/cache"
	"github.com/aity-cloud/monty/pkg/alerting/extensions/destination"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	"github.com/aity-cloud/monty/pkg/logger"
	"log/slog"

	// add profiles
	_ "net/http/pprof"
)

var defaultSeverity = alertingv1.MontySeverity_Info.String()

type EmbeddedServer struct {
	logger *slog.Logger
	// maxSize of the combined caches
	lub int
	// layered caches
	notificationCache cache.MessageCache[alertingv1.MontySeverity, *alertingv1.MessageInstance]
	alarmCache        cache.MessageCache[alertingv1.MontySeverity, *alertingv1.MessageInstance]
	sendK8s           bool
	k8sDestination    destination.Destination
}

func NewEmbeddedServer(
	lg *slog.Logger,
	lub int,
	sendK8s bool,
) *EmbeddedServer {
	e := &EmbeddedServer{
		logger:            lg,
		sendK8s:           sendK8s,
		lub:               lub,
		notificationCache: cache.NewLFUMessageCache(lub),
		alarmCache:        cache.NewLFUMessageCache(lub),
	}
	if sendK8s {
		e.logger.Info("Configuring alerts to be sent to kubernetes events...")
		k8s, err := destination.NewK8sDestination(lg)
		if err != nil {
			panic(err)
		}
		e.k8sDestination = k8s

	}
	return e
}

func StartMontyEmbeddedServer(
	ctx context.Context,
	montyAddr string,
	sendK8s bool,
) *http.Server {
	lg := logger.NewPluginLogger().WithGroup("monty.alerting")
	es := NewEmbeddedServer(lg, 125, sendK8s)
	mux := http.NewServeMux()

	// request body will be in the form of AM webhook payload :
	// https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
	//
	// Note :
	//    Webhooks are assumed to respond with 2xx response codes on a successful
	//	  request and 5xx response codes are assumed to be recoverable.
	// therefore, non-recoverable errors should have error codes 3XX and 4XX
	mux.HandleFunc(shared.AlertingDefaultHookName, es.handleWebhook)
	mux.HandleFunc("/notifications/list", es.handleListNotifications)
	mux.HandleFunc("/alarms/list", es.handleListAlarms)

	hookServer := &http.Server{
		// explicitly set this to 0.0.0.0 for test environment
		Addr:    montyAddr,
		Handler: mux,
	}
	go func() {
		lg.With("addr", montyAddr).Info("starting monty embedded server")
		err := hookServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	go func() {
		<-ctx.Done()
		if err := hookServer.Close(); err != nil {
			panic(err)
		}
	}()
	return hookServer
}
