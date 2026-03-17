package config

type AppConfig struct {
	MySQL   MySQLConfig
	Redis   RedisConfig
	Jwt     JwtConfig
	RunPort string
}

type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JwtConfig struct {
	Secret string
}
