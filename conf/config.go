package conf

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
)

const (
	envInitConfig     = "INIT_CONFIG"
	defaultInitConfig = "environment"

	// 默认给 default serviceaccount 配置 imagePullSecrets
	defautServiceAccont = "default"
	// 默认处理所有 ns
	defautWatchNamespace = "*"

	secretNamePrefix = "image-credential-helper-"
	emailSuffix      = "@imagecredentialhelper.io"

	defaultHttpHealthCheckPort = "8080"
)

type Config struct {
	initConfigFrom      string
	ForceUpdateSecret   bool
	HttpHealthCheckPort string
	*ImageCredentialInfo
}

type ImageCredentialInfo struct {
	// 存储 ImageCredential 的 SecretName，将默认加上前缀 secretNamePrefix
	SecretName string

	Host     string
	User     string
	Password string
	Email    string

	ServiceAccounts []string
	WatchNamespaces []string
}

func (c Config) print() {
	glog.Infof("InitConfigFrom: %s", c.initConfigFrom)
	glog.Infof("ForceUpdateSecret: %t", c.ForceUpdateSecret)
	glog.Infof("HttpHealthCheckPort: %s", c.HttpHealthCheckPort)
	glog.Infof("ImageCredentialInfo.SecretName: %s", c.ImageCredentialInfo.SecretName)
	glog.Infof("ImageCredentialInfo.Host: %s", c.ImageCredentialInfo.Host)
	glog.Infof("ImageCredentialInfo.User: %s", c.ImageCredentialInfo.User)
	glog.Infof("ImageCredentialInfo.Email: %s", c.ImageCredentialInfo.Email)
	glog.Infof("ImageCredentialInfo.ServiceAccounts: %+v", c.ImageCredentialInfo.ServiceAccounts)
	glog.Infof("ImageCredentialInfo.WatchNamespaces: %+v", c.ImageCredentialInfo.WatchNamespaces)
}

func InitConfig() (*Config, error) {
	var err error
	ic := strings.TrimSpace(os.Getenv(envInitConfig))
	if ic == "" {
		return nil, fmt.Errorf("please define env %s", envInitConfig)
	}

	var ici *Config

	switch ic {
	case defaultInitConfig:
		envLoader := EnvLoader{}
		ici, err = envLoader.Load()
		if err != nil {
			return nil, err
		}
	}
	if ici == nil {
		return nil, fmt.Errorf("config is nil")
	}

	ici.print()
	return ici, nil
}
