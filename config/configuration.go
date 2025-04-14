package config

import (
	"flag"
)

type ServerName string

type DevNodeConfig struct {
	Port string
}

type ServerConfig struct {
	Name          ServerName
	Port          string
	DevNodeConfig DevNodeConfig
}

func GetServerConfigFromFlag(name ServerName) ServerConfig {
	var port string
	var serverPort string
	flag.StringVar(&port, "port", "8545", "HTTP RPC port for the dev node")
	flag.StringVar(&serverPort, "serverPort", "8888", "HTTP RPC port for the supporting server")
	flag.Parse()

	return ServerConfig{
		Name: name,
		Port: serverPort,
		DevNodeConfig: DevNodeConfig{
			Port: port,
		},
	}
}

func GetServerConfig(name ServerName) ServerConfig {
	devNodeConfig := DevNodeConfig{
		Port: "8565",
	}
	registry := make(map[ServerName]ServerConfig)
	registry["DevServer"] = ServerConfig{
		Name:          "DevServer",
		Port:          "8575",
		DevNodeConfig: devNodeConfig,
	}
	registry["LogServer"] = ServerConfig{
		Name:          "LogServer",
		Port:          "9585",
		DevNodeConfig: devNodeConfig,
	}

	return registry[name]
}

func (config ServerConfig) GetServerUrl(pathSegment string) string {
	return "http://localhost:" + config.Port + "/" + pathSegment
}
