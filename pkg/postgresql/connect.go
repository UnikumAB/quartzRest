package postgresql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func ConnectPostgresql(dataSourceName string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %s", err)
	}
	logrus.Debug("Connected to Postgresql database")
	return db
}
