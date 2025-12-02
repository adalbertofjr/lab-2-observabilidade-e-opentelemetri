package configs

import "github.com/spf13/viper"

type Conf struct {
	WeatherAPIKey string `mapstructure:"WEATHERAPI_KEY"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Bind explicitamente as variáveis de ambiente
	viper.BindEnv("WEATHERAPI_KEY")
	viper.BindEnv("WEB_SERVER_PORT")

	// Tenta ler .env, mas ignora se não existir
	_ = viper.ReadInConfig()

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
