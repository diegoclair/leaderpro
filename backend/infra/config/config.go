package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config      *Config
	configError error
	once        sync.Once
)

// EnvKeyReplacer replace for environment variable parse
var EnvKeyReplacer = strings.NewReplacer(".", "_", "-", "_")

func setup() {
	var env string
	env = os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	// Get the directory of the current file (config.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		slog.Error("Unable to get the current file path")
		return
	}

	configFilePath := filepath.Join(filepath.Dir(filename),
		"..", "..", "deployment",
		fmt.Sprintf("config-%s.toml", env),
	)
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()
}

// GetConfigEnvironment read config from environment variables and config.toml file
func GetConfigEnvironment(ctx context.Context, appName string) (*Config, error) {
	once.Do(func() {

		setup()

		configError = viper.ReadInConfig()
		if configError != nil {
			slog.Error("Error to read configs: ", slog.String("error", configError.Error()))
			return
		}

		// set default variables on viper
		// for example, if we have a file .toml like this:
		// [db]
		// pass
		// then the ley will be DB_PASS and if we find this key on environment variables, we will set the value on viper
		for _, k := range viper.AllKeys() {
			key := strings.ToUpper(EnvKeyReplacer.Replace(k))
			envValue := os.Getenv(key)
			if envValue != "" {
				viper.Set(k, envValue) // set as default (ignoring config file value)
			}
		}

		config = &Config{
			ctx:     ctx,
			appName: appName,
		}
		configError = viper.Unmarshal(config)
		if configError != nil {
			slog.Error("Error to unmarshal configs: ", slog.String("error", configError.Error()))
			return
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			if in.Op == fsnotify.Write {
				err := viper.Unmarshal(config)
				if err != nil {
					slog.Error("Error to unmarshal new config changes: ", slog.String("error", err.Error()))
					return
				}
			}
		})

		config.setupTracer()
	})

	return config, configError
}
