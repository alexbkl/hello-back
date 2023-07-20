package config

type Config struct {
	// App env
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`
	// Postgres
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	// Filebase credential
	FilebaseAccessKey  string `mapstructure:"FILEBASE_ACCESS_KEY"`
	FilebaseSecretKey  string `mapstructure:"FILEBASE_SECRET_KEY"`
	FilebasePinningKey string `mapstructure:"FILEBASE_PINNING_KEY"`
}
