package model

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/UnikumAB/quartzRest/pkg/types"
)

type Triggers struct {
	SchedName    string           `db:"sched_name"`
	TriggerName  string           `db:"trigger_name"`
	TriggerGroup string           `db:"trigger_group"`
	JobName      string           `db:"job_name"`
	JobGroup     string           `db:"job_group"`
	Description  sql.NullString   `db:"description"`
	NextFireTime types.NanoTime   `db:"next_fire_time"`
	PrevFireTime types.NanoTime   `db:"prev_fire_time"`
	Priority     int              `db:"priority"`
	TriggerState TriggerStateEnum `db:"trigger_state"`
	TriggerType  string           `db:"trigger_type"`
	StartTime    types.NanoTime   `db:"start_time"`
	EndTime      types.NanoTime   `db:"end_time"`
	CalendarName sql.NullString   `db:"calendar_name"`
	MisfireInstr int              `db:"misfire_instr"`
	JobData      []byte           `db:"job_data"`
}

type TriggerStateEnum string

const (
	ERROR    TriggerStateEnum = "ERROR"
	WAITING  TriggerStateEnum = "WAITING"
	BLOCKED  TriggerStateEnum = "BLOCKED"
	ACQUIRED TriggerStateEnum = "ACQUIRED"
)

func (t *TriggerStateEnum) Scan(src interface{}) error {
	if src == nil {
		return errors.New("cannot map null value")
	}
	if src, ok := src.(string); ok {
		switch src {
		case string(ERROR):
			*t = ERROR
			return nil
		case string(WAITING):
			*t = WAITING
			return nil
		case string(BLOCKED):
			*t = BLOCKED
			return nil
		case string(ACQUIRED):
			*t = ACQUIRED
			return nil
		}
		return errors.New("failed to parse string " + src)
	}
	return fmt.Errorf("cannot parse type %T with value %v", src, src)
}

func (t TriggerStateEnum) Value() (driver.Value, error) {
	return string(t), nil
}
