package conf

import (
	"fmt"
	"os"
	"strings"
)

const (
	// 私有仓库产品名，会体现在 secret docker-registry 名字中
	// 比如 harbor
	imageProvider = "IMAGE_PROVIDER"

	envHost     = "IMAGE_HOST"
	envUser     = "IMAGE_USER"
	envPassword = "IMAGE_PASSWORD"

	envServiceAccounts = "SERVICE_ACCOUNTS"
	envWatchNamespaces = "WATCH_NAMESPACES"
)

var _ IImageCredentialInfoLoader = new(EnvImageCredentialInfoLoader)

type EnvImageCredentialInfoLoader struct{}

func (e EnvImageCredentialInfoLoader) Load() (*ImageCredentialInfo, error) {
	ip := strings.ToLower(strings.TrimSpace(os.Getenv(imageProvider)))
	if ip == "" {
		return nil, fmt.Errorf("please define env %s", imageProvider)
	}

	host := strings.TrimSpace(os.Getenv(envHost))
	if host == "" {
		return nil, fmt.Errorf("please define env %s", envHost)
	}

	user := strings.TrimSpace(os.Getenv(envUser))
	if user == "" {
		return nil, fmt.Errorf("please define env %s", envUser)
	}

	pw := strings.TrimSpace(os.Getenv(envPassword))
	if pw == "" {
		return nil, fmt.Errorf("please define env %s", envPassword)
	}

	sas := strings.TrimSpace(os.Getenv(envServiceAccounts))
	if sas == "" {
		sas = defautServiceAccont
	}
	wns := strings.TrimSpace(os.Getenv(envWatchNamespaces))
	if wns == "" {
		wns = defautWatchNamespace
	}

	sasSps := strings.Split(sas, ",")
	wnsSps := strings.Split(wns, ",")
	for _, ns := range wnsSps {
		if ns == "*" {
			wnsSps = []string{defautWatchNamespace}
			break
		}
	}

	return &ImageCredentialInfo{
		SecretName:      secretNamePrefix + ip,
		Host:            host,
		User:            user,
		Password:        envPassword,
		Email:           user + emailSuffix,
		ServiceAccounts: sasSps,
		WatchNamespaces: wnsSps,
	}, nil
}
