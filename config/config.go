package config

type (
	Config struct {
		DBMongo DBMongo
		DBRedis DBRedis
		Log     Log
		Secrets Secrets
		Server  Server
	}

	DBMongo struct {
		URI        string
		DBName     string
		Collection string
		Timeout    int
	}

	DBRedis struct {
		Addr string
		Port string
		User string
		Pass string
	}

	Log struct {
		FilePath string
	}
	Server struct {
		Addr string
		Port string
	}

	Secrets struct {
		JwtSecret string `envconfig:"JWT_SECRET"`
		ExpTime   int    `envconfig:"EXP_TIME"`
	}
)
