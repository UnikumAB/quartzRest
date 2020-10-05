package model

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/UnikumAB/quartzRest/pkg/types"
)

type Triggers struct {
	SchedName    string
	TriggerName  string
	TriggerGroup string
	JobName      string
	JobGroup     string
	Description  string
	NextFireTime types.NanoTime
	PrevFireTime types.NanoTime
	Priority     int
	TriggerState TriggerStateEnum
	TriggerType  string
	StartTime    types.NanoTime
	EndTime      types.NanoTime
	CalendarName string
	MisfireInstr int
	JobData      []byte
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
