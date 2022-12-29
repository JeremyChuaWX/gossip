package initialisers

import "github.com/spf13/viper"

type Config struct {
	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBName     string `mapstructure:"POSTGRES_DBNAME"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`

	ServerPort string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
