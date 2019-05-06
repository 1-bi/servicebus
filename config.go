package servicebus

import "github.com/coreos/etcd/clientv3"

// Config
type Config struct {
	etcdConfig *clientv3.Config
}

func NewConfig() *Config {
	var conf = new(Config)
	return conf
}

// SetRegisterServer define restier servers handle
func (myself *Config) SetRegisterServer(config *clientv3.Config) {
	myself.etcdConfig = config

}

func (myself *Config) SetSettingDir(paths ...string) {

}
