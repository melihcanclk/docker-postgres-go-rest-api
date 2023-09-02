package config

import (
	"os"
	"strconv"
)

var AccessTokenPrivateKey string = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
var AccessTokenPublicKey string = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
var AccessTokenExpiredInMinutes = os.Getenv("ACCESS_TOKEN_EXPIRED_IN")
var AccessTokenMaxAge, _ = strconv.ParseInt(os.Getenv("ACCESS_TOKEN_MAX_AGE"), 10, 0)

var RefreshTokenPrivateKey string = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
var RefreshTokenPublicKey string = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
var RefreshTokenExpiredInMinutes = os.Getenv("REFRESH_TOKEN_EXPIRED_IN")
var RefreshTokenMaxAge, _ = strconv.ParseInt(os.Getenv("REFRESH_TOKEN_MAX_AGE"), 10, 0)
