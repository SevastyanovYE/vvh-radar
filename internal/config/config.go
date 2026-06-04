package config

import (
	"os"
)

type Config struct {
	DBPath  string
	Addr    string
	GroupID string
}

func Load() Config {
	return Config{
		DBPath:  getEnv("VVH_RADAR_DB", "./vvh-radar.db"),
		Addr:    getEnv("VVH_RADAR_ADDR", ":8080"),
		GroupID: getEnv("VVH_RADAR_GROUP_ID", "fake-group"),
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
