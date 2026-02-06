package config

type Config struct {
	AppPort string
	DBURL   string
}

func Load() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DBURL:   getEnv("DATABASE_URL", ""),
	}
}
