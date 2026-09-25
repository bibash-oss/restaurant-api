package config

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config interface {
	IsSet(key string) bool
	GetString(key string) string
	GetInt(key string) int
	GetFloat(key string) float64
	GetBool(key string) bool
	GetIntSlice(key string) []int
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
}

type ViperConfig struct{}

var (
	once     sync.Once
	instance Config
)

func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

func initConfigs() {
	loadDotEnv(".env")

	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if configFilePath := os.Getenv("CONFIG_FILE_PATH"); configFilePath != "" {
		stat, err := os.Stat(configFilePath)
		if err == nil {
			if stat.IsDir() {
				viper.AddConfigPath(configFilePath)
			} else {
				viper.SetConfigFile(configFilePath)
			}
		}
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: error reading config file: %v", err)
	}
}

func Default() Config {
	once.Do(func() {
		initConfigs()
		instance = &ViperConfig{}
	})
	return instance
}

func (v *ViperConfig) IsSet(key string) bool {
	return viper.IsSet(key)
}

func (v *ViperConfig) GetString(key string) string {
	return viper.GetString(key)
}

func (v *ViperConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

func (v *ViperConfig) GetFloat(key string) float64 {
	return viper.GetFloat64(key)
}

func (v *ViperConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (v *ViperConfig) GetIntSlice(key string) []int {
	return viper.GetIntSlice(key)
}

func (v *ViperConfig) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (v *ViperConfig) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func (v *ViperConfig) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}
