package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davidhong1/k8s-image-credential-helper/conf"
	"github.com/davidhong1/k8s-image-credential-helper/service"
	"github.com/golang/glog"
)

func main() {
	ctx := context.Background()

	config, err := conf.InitConfig()
	if err != nil {
		glog.Fatal("Init config fail.", err)
	}

	nsWatcher, err := service.InitNamespaceWatcher(ctx, config)
	if err != nil {
		glog.Fatal(err)
	}
	go func(ctx context.Context) {
		startWatchNum := 1
		for {
			if startWatchNum > 1 {
				glog.Info("watcher.ResultChan is closed, start new watcher, startWatchNum: %d", startWatchNum)
			}
			err := nsWatcher.Watch(ctx)
			if err != nil {
				glog.Fatal(err)
			}
			nsWatcher.ForceUpdateSecret = false
			startWatchNum++
		}
	}(ctx)

	// add http health check api /pong
	http.HandleFunc("/pong", pongHandler)
	err = http.ListenAndServe(":"+config.HttpHealthCheckPort, nil)
	if err != nil {
		glog.Fatal("Init health http failed. err: %v", err)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "ping")
}
