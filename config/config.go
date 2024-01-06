package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	Port            string
	ReadTimeOut     int
	WriteTimeOut    int
	RequestTimeOut  int
	ShutdownTimeOut int

	PostgresAddresses []string
	PostgresDBName    string
	PostgresUser      string
	PostgresPassword  string

	TDXClientID     string
	TDXClientSecret string
)

func Initialize(path string) {
	Port = ":25976"
	ReadTimeOut = 180
	WriteTimeOut = 60
	RequestTimeOut = 60
	ShutdownTimeOut = 15

	PostgresAddresses = []string{"127.0.0.1:5432"}
	PostgresDBName = "bus"
	PostgresUser = "postgres"
	PostgresPassword = "postgres"

	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	viper.AddConfigPath(".")

	viper.SetDefault("PORT", ":12345")
	viper.SetDefault("READ_TIMEOUT", 30)
	viper.SetDefault("WRITE_TIMEOUT", 30)
	viper.SetDefault("REQUEST_TIMEOUT", 30)
	viper.SetDefault("SHUTDOWN_TIMEOUT", 10)

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	Port = viper.GetString("PORT")
	ReadTimeOut = viper.GetInt("READ_TIMEOUT")
	WriteTimeOut = viper.GetInt("WRITE_TIMEOUT")
	RequestTimeOut = viper.GetInt("REQUEST_TIMEOUT")
	ShutdownTimeOut = viper.GetInt("SHUTDOWN_TIMEOUT")

	PostgresAddresses = viper.GetStringSlice("POSTGRES_ADDRESSES")
	PostgresDBName = viper.GetString("POSTGRES_DB_NAME")
	PostgresUser = viper.GetString("POSTGRES_USER")
	PostgresPassword = viper.GetString("POSTGRES_PASSWORD")

	TDXClientID = viper.GetString("TDX_CLIENT_ID")
	TDXClientSecret = viper.GetString("TDX_CLIENT_SECRET")
}
