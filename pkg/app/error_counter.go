package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

func (app App) updateErrorCounter(errorCounter *prometheus.CounterVec) {

	query := "SELECT sched_name, trigger_group, trigger_name, count(*) " +
		"FROM " + app.Prefix + "triggers " +
		"WHERE trigger_state='ERROR' " +
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
		count := 0
		err := result.Scan(&schedName, &triggerGroup, &triggerName, &count)
		if err != nil {
			logrus.Errorf("scanning result failed: %v", err)
		}
		labels := prometheus.Labels{
			"sched_name":    schedName,
			"trigger_group": triggerGroup,
			"trigger_name":  triggerName,
		}
		errorCounter.With(labels).Inc()
	}
}

func (app App) createErrorCounter() *prometheus.CounterVec {
	errorCounter := promauto.NewCounterVec(prometheus.CounterOpts{Name: "quartz_error_counter", Help: "No Help"},
		[]string{"sched_name", "trigger_group", "trigger_name"})
	return errorCounter
}
