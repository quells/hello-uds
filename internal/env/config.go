package env

import "github.com/kelseyhightower/envconfig"

var Config struct {
	SocketFile string `envconfig:"SOCKET_FILE" default:"/tmp/hello-uds.sock"`
}

func init() {
	envconfig.MustProcess("", &Config)
}
