// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Jwt   JwtConfig
	Mysql MysqlConf
	Redis RedisConf
}

type JwtConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpire  int64
	RefreshTokenExpire int64
}

type MysqlConf struct {
	Dsn             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConf struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}
