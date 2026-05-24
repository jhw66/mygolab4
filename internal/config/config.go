// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Jwt   JwtConfig
	Totp  TotpConfig
	Mysql MysqlConf
	Redis RedisConf
}

type TotpConfig struct {
	Issuer          string
	SecretCipherKey string
	QRCodeSize      int
	QRCodeLevel     string
}

type JwtConfig struct {
	AccessTokenSecret    string
	RefreshTokenSecret   string
	ChallengeTokenSecret string
	AccessTokenExpire    int64
	RefreshTokenExpire   int64
	ChallengeTokenExpire int64
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
