package config

import "os"

func devMode() bool {
	return os.Getenv("DEV") == "1"
}

func GetDbUrl() string {
	if devMode() {
		return os.Getenv("DEV_DBURL")
	} else {
		return os.Getenv("DBURL")
	}
}
