package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	C    = new(Config)
	once sync.Once
)

// Load config file (toml/json/yaml)
func MustLoad(path string) {
	once.Do(func() {
		viper.SetConfigFile(path)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = viper.Unmarshal(&C)
		if err != nil {
			panic(err)
		}

	})
}

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

type Config struct {
	RunMode     string    `json:"runMode"`
	PrintConfig bool      `json:"printConfig"`
	HTTP        HTTP      `json:"http"`
	Commands    []Command `json:"commands"`
}

type HTTP struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Method string `json:"method"`
}

type Command struct {
	Code string `json:"code"`
	Exec string `json:"exec"`
}
