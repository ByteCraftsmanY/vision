package config

import (
	"flag"
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"strings"
	"vision/constants/enums"
)

type Configurations struct {
	Environment string
	Server      ServerConfigurations
	Database    DatabaseConfigurations
	JWT         JWTConfigurations
	Kafka       KafkaConfigurations
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

type KafkaConfigurations struct {
	DefaultTopic string
	Brokers      []string
}

func (c Configurations) IsProductionEnvironment() bool {
	return strings.EqualFold(c.Environment, enums.EnvProduction.String())
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
