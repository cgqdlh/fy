package fy

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
)

//go:embed config-example.yml
var configExample []byte
var (
	once         sync.Once
	cfg          Config
	callpointBuf strings.Builder
)

func GetConfig() Config {
	once.Do(func() {
		cfg := initConfig()
		if level, err := log.ParseLevel(cfg.Log); err != nil {
			initLog(log.WarnLevel)
		} else {
			initLog(level)
		}

		cfgJson, err := json.Marshal(&cfg)
		if err != nil {
			log.Fatalf("Fatal configuration converted to json: %s\n", err)
		}
		pretty.DefaultOptions.SortKeys = true
		res := pretty.Color(pretty.PrettyOptions(cfgJson, pretty.DefaultOptions), nil)
		log.Debugf(":\n%s", res)
	})
	return cfg
}

func initConfig() Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("$GO_FY_CONFIG_FILE")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/go-fy")
	viper.AddConfigPath("/etc/go-fy")
	viper.SetEnvPrefix("go_fy")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			msg := fmt.Sprintf("No configuration file found: %s\nexample: config.yml\n```\n%s\n```\n", err, configExample)
			log.Fatalf(msg)
		} else {
			log.Fatalf("Configuration file error: %s\n", err)
		}
	}
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("configuration file parsing error: %s\n", err)
	}
	return cfg
}

func defaultConfig() {
	viper.SetDefault("log", "warn")
}

func initLog(level log.Level) {
	formatter := &log.TextFormatter{
		CallerPrettyfier: callerPrettyfier,
		TimestampFormat:  "2006-01-02T15:04:05.000Z07:00",
		FullTimestamp:    true,
	}
	log.SetFormatter(formatter)
	log.SetLevel(level)
	if level >= log.DebugLevel {
		log.SetReportCaller(true)
	}
}

func callerPrettyfier(f *runtime.Frame) (function string, file string) {
	callpointBuf.Reset()
	callpointBuf.WriteByte('(')
	callpointBuf.WriteString(f.Function)
	callpointBuf.WriteByte(':')
	callpointBuf.WriteString(strconv.Itoa(f.Line))
	callpointBuf.WriteByte(')')
	return callpointBuf.String(), ""
}
