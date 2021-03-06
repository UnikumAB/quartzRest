package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

func (app App) updateErrorCounter(errorCounter *prometheus.CounterVec,
	blockedCounter *prometheus.CounterVec) {
	query := "SELECT sched_name, trigger_group, trigger_name, trigger_state, count(*) " +
		"FROM " + app.Prefix + "triggers " +
		"WHERE trigger_state in ('ERROR','BLOCKED') " +
		"GROUP BY sched_name, trigger_group, trigger_name"
	result, err := app.DB.Query(query)

	if err != nil {
		logrus.Errorf("Querying for errored triggers failed: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer result.Close()

	for result.Next() {
		schedName := ""
		triggerName := ""
		triggerGroup := ""
		triggerState := ""
		count := 0
		err := result.Scan(&schedName, &triggerGroup, &triggerName, &triggerState, &count)

		if err != nil {
			logrus.Errorf("scanning result failed: %v", err)
		}

		labels := prometheus.Labels{
			"sched_name":    schedName,
			"trigger_group": triggerGroup,
			"trigger_name":  triggerName,
		}

		switch triggerState {
		case "ERROR":
			errorCounter.With(labels).Inc()
		case "BLOCKED":
			blockedCounter.With(labels).Inc()
		}
	}
}

func (app App) createErrorCounter() *prometheus.CounterVec {
	errorCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "quartz_error_counter",
		Help: "Number of Jobs encountered that are in ERROR state"},
		[]string{"sched_name", "trigger_group", "trigger_name"})

	return errorCounter
}
func (app App) createBlockedCounter() *prometheus.CounterVec {
	errorCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "quartz_blocked_counter",
		Help: "Number of Jobs encountered that are in BLOCKED state"},
		[]string{"sched_name", "trigger_group", "trigger_name"})

	return errorCounter
}
