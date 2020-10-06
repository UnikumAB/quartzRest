package main

import (
	"os"

	"github.com/UnikumAB/quartzRest/pkg/app"
	"github.com/UnikumAB/quartzRest/pkg/postgresql"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	a := app.App{}
	pgConn := ""
	cmd := kingpin.New(os.Args[0], "The webserver")
	cmd.Flag("port", "the port to run on").
		Default("localhost:8080").
		Envar("QUARTZ_SERVER_PORT").
		StringVar(&a.Port)
	cmd.Flag("postgres-connection", "Connection string for the postgres database").
		Short('P').
		Envar("QUARTZ_SERVER_POSTGRESS_CONNECTION").
		StringVar(&pgConn)
	cmd.Flag("table-prefix", "Prefix of the quartz tables").
		Default("qrtz_").
		Envar("QUARTZ_SERVER_PREFIX").
		StringVar(&a.Prefix)
	kingpin.MustParse(cmd.Parse(os.Args[1:]))
	if pgConn != "" {
		a.DB = postgresql.ConnectPostgresql(pgConn)
		err := a.DB.Ping()
		if err != nil {
			logrus.Fatalf("Failed to ping database: %v", err)
		}
		logrus.Infof("using database %v (%v)", pgConn, a.DB.DriverName())
	}
	if a.DB == nil {
		logrus.Fatal("Need to provide some database connection.")
	}
	if a.Port == "" {
		logrus.Fatal("Port may not be null")
	}
	a.Run()
}
