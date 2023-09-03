package config

import "os"

var RedisPort = os.Getenv("REDIS_URL")
