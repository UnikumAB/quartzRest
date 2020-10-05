package model

import (
	"github.com/UnikumAB/quartzRest/pkg/types"
)

type SchedulerState struct {
	SchedName       string         `db:"sched_name"`
	InstanceName    string         `db:"instance_name"`
	LastCheckinTime types.NanoTime `db:"last_checkin_time"`
	CheckinInterval int            `db:"checkin_interval"`
}
