package config

type Config struct {
	MongoConfig MongoConfig
}

type MongoConfig struct {
	Host string
}

//
func Load() Config {
	return loadDefault()
}

//
func loadDefault() Config {
	cfg := Config{
		MongoConfig: MongoConfig{
			Host: "localhost",
		},
	}

	return cfg
}
