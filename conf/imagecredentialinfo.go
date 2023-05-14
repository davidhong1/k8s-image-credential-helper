package conf

const (
	// 默认给 default serviceaccount 配置 imagePullSecrets
	defautServiceAccont = "default"
	// 默认处理所有 ns
	defautWatchNamespace = "*"

	secretNamePrefix = "image-credential-helper-"
	emailSuffix      = "@imagecredentialhelper.io"
)

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

type IImageCredentialInfoLoader interface {
	Load() (*ImageCredentialInfo, error)
}
