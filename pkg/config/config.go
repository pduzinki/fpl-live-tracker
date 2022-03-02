package config

type Config struct {
	MongoConfig MongoConfig
}

type MongoConfig struct {
	Host     string
	Port     int
	Database string
}

//
func Load() Config {
	return loadDefault()
}

//
func loadDefault() Config {
	cfg := Config{
		MongoConfig: MongoConfig{
			Host:     "localhost",
			Port:     27017,
			Database: "fpl-live-tracker",
		},
	}

	return cfg
}
