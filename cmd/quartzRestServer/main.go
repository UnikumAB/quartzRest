package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/UnikumAB/quartzRest/pkg/app"
	"github.com/UnikumAB/quartzRest/pkg/postgresql"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "QuartzRest is a simple tool for accessing quartz data tables.",
		Run: func(cmd *cobra.Command, args []string) {
			a := app.App{}

			url, err := url.Parse(viper.GetString("postgres-connection"))
			if err != nil {
				logrus.Fatalf("Failed parse database url: %v", err)
			}
			if url.Scheme == "postgres" {
				a.DB = postgresql.ConnectPostgresql(url.String())
			}
			err = a.DB.Ping()
			if err != nil {
				logrus.Fatalf("Failed to ping database: %v", err)
			}

			logrus.Infof("Connected to DB %q", url.Redacted())
			a.Port = viper.GetString("bind")
			logrus.Infof("Listening on %q", a.Port)
			a.Prefix = viper.GetString("table-prefix")
			logrus.Infof("Using table prefix %q", a.Prefix)

			a.Run()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("bind", "b", "localhost:8080", "the host:port to bind to")
	rootCmd.PersistentFlags().StringP("postgres-connection", "P", "", "Connection string for the postgres database")
	rootCmd.PersistentFlags().StringP("table-prefix", "p", "qrtz_", "Prefix of the quartz tables")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	mustComplete(viper.BindPFlag("bind", rootCmd.PersistentFlags().Lookup("bind")))
	mustComplete(viper.BindPFlag("postgres-connection", rootCmd.PersistentFlags().Lookup("postgres-connection")))
	mustComplete(viper.BindPFlag("table-prefix", rootCmd.PersistentFlags().Lookup("table-prefix")))
	mustComplete(viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper")))
	viper.SetDefault("bind", "localhost:8080")
	viper.SetDefault("table-prefix", "qrtz_")

	viper.SetEnvPrefix("QUARTZ_SERVER")
	viper.AutomaticEnv()
}

func mustComplete(err error) {
	if err != nil {
		logrus.Fatalf("Failure: %v", err)
	}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}
