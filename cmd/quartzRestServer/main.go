package main

import (
	"os"

	"github.com/UnikumAB/quartzRest/pkg/app"
	"github.com/UnikumAB/quartzRest/pkg/postgresql"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cmd = kingpin.New(os.Args[0], "The webserver")

	bindto = cmd.Flag("bind", "the host:port to bind to").
		Default("localhost:8080").
		Envar("QUARTZ_SERVER_BIND").
		String()

	pg = cmd.Flag("postgres-connection", "Connection string for the postgres database").
		Short('P').
		Envar("QUARTZ_SERVER_POSTGRESS_CONNECTION").
		Required().
		String()

	prefix = cmd.Flag("table-prefix", "Prefix of the quartz tables").
		Default("qrtz_").
		Envar("QUARTZ_SERVER_PREFIX").
		String()
)

func main() {
	a := app.App{}
	kingpin.MustParse(cmd.Parse(os.Args[1:]))
	a.DB = postgresql.ConnectPostgresql(*pg)
	a.Port = *bindto
	a.Prefix = *prefix

	err := a.DB.Ping()
	if err != nil {
		logrus.Fatalf("Failed to ping database: %v", err)
	}
	a.Run()
}
