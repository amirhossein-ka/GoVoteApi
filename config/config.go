package config

type (
	Config struct {
		DBMongo DBMongo
		DBRedis DBRedis
		Log     Log
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
	}
)
