package config

import legacyConfig "ez-admin-gin/server/internal/config"

type Config = legacyConfig.Config
type AppConfig = legacyConfig.AppConfig
type ServerConfig = legacyConfig.ServerConfig
type AuthConfig = legacyConfig.AuthConfig
type LogConfig = legacyConfig.LogConfig
type DatabaseConfig = legacyConfig.DatabaseConfig
type RedisConfig = legacyConfig.RedisConfig
type UploadConfig = legacyConfig.UploadConfig

// Load 读取应用配置。
func Load() (*Config, error) {
	return legacyConfig.Load()
}
