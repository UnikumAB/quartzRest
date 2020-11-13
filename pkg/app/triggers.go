package app

import (
	"encoding/json"
	"net/http"

	"github.com/UnikumAB/quartzRest/pkg/model"
	"github.com/doug-martin/goqu"
	"github.com/sirupsen/logrus"
)

func (app App) TriggerHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var triggers []model.Triggers

		from := goqu.From(app.Prefix + "triggers")

		//query := "SELECT * FROM " + app.Prefix + "triggers"
		state := request.URL.Query().Get("state")
		expression := goqu.Ex{}

		if state != "" {
			expression["trigger_state"] = state
			expression = goqu.Ex{"trigger_state": state}
		}

		expressions, err := expression.ToExpressions()

		if err == nil && len(expressions.Expressions()) > 0 {
			from = from.Where(expressions)
		}

		sql, args, err := from.ToSql()
		if err != nil {
			logrus.Fatalf("Failed to generate SQL: %s", err)
		}

		err = app.DB.Select(&triggers, sql, args...)

		if err != nil {
			logrus.Fatalf("Failed to select trigger: %v", err)
		}

		if triggers == nil {
			triggers = []model.Triggers{}
		}

		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(triggers)

		if err != nil {
			writer.WriteHeader(500)
			logrus.Errorf("Failed to handle request")
		}
	}
}
