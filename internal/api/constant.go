package api

import "github.com/spf13/viper"

const DBDriverName = "postgres"

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("IS_HTTP_LOGGER", false)
	viper.SetDefault("CORS_WILDCARD", "*") // or set current domain for CORS

	viper.SetDefault("LOGSTASH_HOST", "localhost")
	viper.SetDefault("LOGSTASH_PORT", 31130)

	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_NAME", "payroll")
	viper.SetDefault("DB_USER", "payroll_db_user")
	viper.SetDefault("DB_PASSWORD", "jw8s0F4")
	viper.SetDefault("DB_SCHEMA", "payroll")
	viper.SetDefault("SSL_MODE", false)

}
