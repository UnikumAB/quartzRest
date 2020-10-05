package types

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type NanoTime struct {
	time.Time
}

func (i *NanoTime) Scan(src interface{}) error {
	switch src := src.(type) {
	case int:
		unix := time.Unix(0, int64(src)*1000*1000)
		i.Time = unix
		logrus.Infof("src=%v, time=%v", src, unix)
		return nil
	case int32:
		unix := time.Unix(0, int64(src)*1000*1000)
		i.Time = unix
		logrus.Infof("src=%v, time=%v", src, unix)
		return nil
	case int64:
		unix := time.Unix(0, src*1000*1000)
		i.Time = unix
		logrus.Infof("src=%v, time=%v", src, unix)
		return nil
	}
	return fmt.Errorf("cannot convert %s (%T) to Time", src, src)
}

func (i NanoTime) Value() (driver.Value, error) {
	return (i.Time.Unix() * 1000) + int64(i.Time.Nanosecond()), nil
}
