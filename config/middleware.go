package config

import "os"

var Secret = os.Getenv("JWT_SECRET")
