package postgresql

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func ConnectPostgresql(dataSourceName string) *sqlx.DB {
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %s", err)
	}

	logrus.Debug("Connected to Postgresql database")

	return db
}
