package config

type (
	Config struct {
		DBRedis    DBRedis
		DBPostgres DBPostgres
		Log        Log
		Secrets    Secrets
		Server     Server
	}

	DBRedis struct {
		Addr string
		Port string
		User string
		Pass string
	}

	DBPostgres struct {
		URI     string
		DBName  string
		Timeout uint16
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
