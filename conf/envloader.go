package conf

import (
	"fmt"
	"os"
	"strings"
)

const (
	envHTTPHealthCheckPort = "HTTP_HEALTH_CHECK_PORT"

	// 私有仓库产品名，会体现在 secret docker-registry 名字中
	// 比如 harbor
	envImageProvider = "IMAGE_PROVIDER"

	envHost     = "IMAGE_HOST"
	envUser     = "IMAGE_USER"
	envPassword = "IMAGE_PASSWORD"

	envServiceAccounts = "SERVICE_ACCOUNTS"
	envWatchNamespaces = "WATCH_NAMESPACES"

	envForceUpdateSecret = "FORCE_UPDATE_SECRET"
)

var _ IConfigLoader = new(EnvLoader)

type EnvLoader struct{}

func (e EnvLoader) Load() (*Config, error) {
	httpHealthCheckPort := strings.ToLower(strings.TrimSpace(os.Getenv(envHTTPHealthCheckPort)))
	if httpHealthCheckPort == "" {
		httpHealthCheckPort = defaultHttpHealthCheckPort
	}

	iProvider := strings.ToLower(strings.TrimSpace(os.Getenv(envImageProvider)))
	if iProvider == "" {
		return nil, fmt.Errorf("please define env %s", envImageProvider)
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

	forceUpdateSecret := false
	fus := strings.ToLower(strings.TrimSpace(os.Getenv(envForceUpdateSecret)))
	if fus == "1" || fus == "yes" || fus == "y" {
		forceUpdateSecret = true
	}

	sasSps := strings.Split(sas, ",")
	wnsSps := strings.Split(wns, ",")
	for _, ns := range wnsSps {
		if ns == "*" {
			wnsSps = []string{defautWatchNamespace}
			break
		}
	}

	return &Config{
		initConfigFrom:      defaultInitConfig,
		ForceUpdateSecret:   forceUpdateSecret,
		HttpHealthCheckPort: httpHealthCheckPort,
		ImageCredentialInfo: &ImageCredentialInfo{SecretName: secretNamePrefix + iProvider, Host: host, User: user, Password: envPassword, Email: user + emailSuffix, ServiceAccounts: sasSps, WatchNamespaces: wnsSps},
	}, nil
}
