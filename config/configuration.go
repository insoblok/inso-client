package config

import "flag"

type DevNodeConfig struct {
	Port string
}

type ServerConfig struct {
	Name          string
	Port          string
	DevNodeConfig DevNodeConfig
}

func GetServerConfigFromFlag(name string) ServerConfig {
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

func GetServerConfig(name string) ServerConfig {
	devNodeConfig := DevNodeConfig{
		Port: "8565",
	}
	registry := make(map[string]ServerConfig)
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
