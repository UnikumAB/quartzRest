package app

import (
	"encoding/json"
	"net/http"

	"github.com/UnikumAB/quartzRest/pkg/model"

	"github.com/sirupsen/logrus"
)

func (app App) SchedulerHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var schedulers []model.SchedulerState
		err := app.DB.Select(&schedulers, "SELECT * FROM "+app.Prefix+"scheduler_state")
		if err != nil {
			logrus.Fatalf("Failed to select schedulers: %v", err)
		}
		if schedulers == nil {
			schedulers = []model.SchedulerState{}
		}
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(schedulers)
		if err != nil {
			writer.WriteHeader(500)
			logrus.Errorf("Failed to handle request")
		}
	}
}
