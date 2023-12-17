package config

import (
	"flag"
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"strings"
	"vision/constants"
)

type Configurations struct {
	Environment string
	Server      ServerConfigurations
	Database    DatabaseConfigurations
	JWT         JWTConfigurations
}

type ServerConfigurations struct {
	Port int
}

type DatabaseConfigurations struct {
	Keyspace string
	Hosts    []string
}

type JWTConfigurations struct {
	Secret string
	Expiry int
}

func (c Configurations) IsProductionEnvironment() bool {
	return strings.EqualFold(c.Environment, constants.EnvProduction.String())
}

var config Configurations

func Init() *Configurations {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	k := koanf.New(".")
	filePath := fmt.Sprintf("config/%s.yml", *environment)
	err := k.Load(file.Provider(filePath), yaml.Parser())
	if err != nil {
		panic(fmt.Sprintf("Failed to load config Error - %v", err))
	}

	config.Environment = *environment
	err = k.Unmarshal("", &config)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal config err - %v", err))
	}
	return &config
}

func GetConfig() *Configurations {
	return &config
}
