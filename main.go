package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/davidhong1/k8s-image-credential-helper/conf"
	"github.com/davidhong1/k8s-image-credential-helper/service"
	"github.com/golang/glog"
)

const (
	httpHealthCheckPort        = "HTTP_HEALTH_CHECK_PORT"
	defaultHttpHealthCheckPort = "8080"

	initConfig        = "INIT_CONFIG"
	defaultInitConfig = "environment"
)

func main() {
	ctx := context.Background()
	var err error

	ic := strings.TrimSpace(os.Getenv(initConfig))
	if ic == "" {
		ic = defaultInitConfig
	}

	var ici *conf.ImageCredentialInfo

	switch ic {
	case defaultInitConfig:
		envLoader := conf.EnvImageCredentialInfoLoader{}
		ici, err = envLoader.Load()
		if err != nil {
			glog.Fatal(err.Error())
		}
	}
	if ici == nil {
		glog.Fatal("ImageCredentialInfo is nil")
	}

	namespaceWatcher, err := service.InitNamespaceWatcher(ctx, ici)
	if err != nil {
		glog.Fatal(err)
	}
	err = namespaceWatcher.Watch(context.Background())
	if err != nil {
		glog.Fatal(err)
	}

	// 监听健康检查接口
	httpPort := strings.TrimSpace(os.Getenv(httpHealthCheckPort))
	if httpPort == "" {
		httpPort = defaultHttpHealthCheckPort
	}
	http.HandleFunc("/pong", pongHandler)
	err = http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		glog.Fatal("init health http failed. err: %v", err)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "ping")
}
