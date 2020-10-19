package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UnikumAB/quartzRest/pkg/promhttpmux"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type App struct {
	DB     *sqlx.DB
	Prefix string
	Port   string
}

func (app App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	r := Router(app)
	srv := &http.Server{
		Addr:    app.Port,
		Handler: handlers.CombinedLoggingHandler(os.Stdout, r),
	}
	schedulerDone := app.runScheduledJobs(ctx)
	shutdownDone := waitForShutdown(cancel, srv, schedulerDone)

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logrus.Fatal(err)
	}
	<-shutdownDone
	logrus.Info("All done. Leaving this cruel world.")
}

// waitForShutdown catches the Interrupt signal and initiates a shutdown.
//
// It will try to shutdown a http.Server if provided, run the provided CancelFunc to cancel the main context and
// then wait for all provided channels to close or send a value. When all is done it will close the returned channel
// to signal that shutdown is complete.
func waitForShutdown(cancel context.CancelFunc, srv *http.Server, channels ...chan struct{}) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		signalChan := make(chan os.Signal, 1)

		signal.Notify(
			signalChan,
			syscall.SIGINT, // kill -SIGINT XXXX or Ctrl+c
		)
		<-signalChan
		cancel()
		timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		if srv != nil {
			err := srv.Shutdown(timeout)
			if err != nil {
				logrus.Errorf("Failed to shutdown server: %v", err)
			}
		}
		for _, channel := range channels {
			<-channel
		}
		logrus.Info("Shutdown complete")
		close(done)
	}()
	return done
}

func Router(app App) *mux.Router {
	r := mux.NewRouter()
	r.Use(promhttpmux.InstrumentHttpDuration(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests."}))
	r.Use(promhttpmux.InstrumentHttpInFlight(prometheus.GaugeOpts{
		Name: "http_in_flight_count",
		Help: "Number of requests in flight"}))
	r.Use(commonMiddleware)
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/schedulers", app.SchedulerHandler())
	r.Handle("/triggers", app.TriggerHandler())
	return r
}

func (app App) runScheduledJobs(ctx context.Context) chan struct{} {
	done := make(chan struct{})
	go func() {
		errorCounter := app.createErrorCounter()
		blockedCounter := app.createBlockedCounter()
		app.updateErrorCounter(errorCounter, blockedCounter)
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				app.updateErrorCounter(errorCounter, blockedCounter)
			case <-ctx.Done():
				logrus.Info("Stopping regular jobs")
				close(done)
				return
			}
		}
	}()
	return done
}
