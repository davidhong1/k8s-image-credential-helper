package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/davidhong1/k8s-image-credential-helper/conf"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type NamespaceWatcher struct {
	clientset *kubernetes.Clientset

	imageCredentialInfo *conf.ImageCredentialInfo

	exitingNamespacesMap sync.Map
}

func InitNamespaceWatcher(ctx context.Context, ici *conf.ImageCredentialInfo) (*NamespaceWatcher, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Error(err)
		return nil, errors.Wrap(err, "")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		return nil, errors.Wrap(err, "")
	}

	return &NamespaceWatcher{
		clientset:            clientset,
		imageCredentialInfo:  ici,
		exitingNamespacesMap: sync.Map{},
	}, nil
}

func (k *NamespaceWatcher) Watch(ctx context.Context) error {
	watcher, err := k.clientset.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		glog.Error(err)
		return errors.Wrap(err, "")
	}

	go func() {
		for event := range watcher.ResultChan() {
			glog.Infof("Watch namespace event type: %v", event.Type)

			switch obj := event.Object.(type) {
			case *v1.Namespace:
				namespace := event.Object.(*v1.Namespace)

				switch event.Type {
				case watch.Added:
					if !k.isWatchNamespace(namespace.Name) {
						glog.Infof("Watch namespace %s created, it not in watch namespaces, will ignore", namespace.Name)
						break
					}

					// 检查是否已处理过，如果已处理，则忽略
					if _, ok := k.exitingNamespacesMap.Load(namespace.Name); ok {
						glog.Infof("Watch namespace %s created, it in exitingNamespacesMap, will ignore", namespace.Name)
						break
					}
					k.exitingNamespacesMap.Store(namespace.Name, namespace.Name)

					glog.Infof("Watch namespace %s created, will create  credential", namespace.Name)
					err := k.createCredential(context.Background(), namespace.Name)
					if err != nil {
						glog.Error(err)
						// TODO 发送通知
					}
				case watch.Deleted:
					glog.Infof("Watch namespace %s deleted", namespace.Name)
					k.exitingNamespacesMap.Delete(namespace.Name)
				}
			default:
				glog.Error("Watch namespace unsupported event object: %T", obj)
			}
		}

		glog.Info("Watch namespace, watcher.ResultChan is closed, start new watcher")
		err := k.Watch(context.Background())
		if err != nil {
			glog.Errorf("Watch namespace, start new watcher failed. err: %v", err)
			// TODO 发送通知
		}
	}()

	return nil
}

func (k *NamespaceWatcher) createCredential(ctx context.Context, ns string) error {
	if !k.isWatchNamespace(ns) {
		return nil
	}

	glog.Infof("createCredential Namespace: %s", ns)

	err := k.createSecret(ctx, ns)
	if err != nil {
		glog.Error(err)
		return errors.Wrap(err, "")
	}

	err = k.updateServiceAccount(ctx, ns)
	if err != nil {
		glog.Error(err)
		return errors.Wrap(err, "")
	}

	return nil
}

func (k *NamespaceWatcher) createSecret(ctx context.Context, ns string) error {
	if !k.isWatchNamespace(ns) {
		return nil
	}

	glog.Infof("createSecret Namespace: %s", ns)

	resp, err := k.clientset.CoreV1().Secrets(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		glog.Error(err)
		return errors.Wrap(err, "")
	}

	for _, v := range resp.Items {
		if v.Name == k.imageCredentialInfo.SecretName {
			// 已存在，直接返回
			return nil
		}
	}

	bs, err := dockerConfigJsonKeyBytes(k.imageCredentialInfo.Host, k.imageCredentialInfo.User, k.imageCredentialInfo.Password, k.imageCredentialInfo.Email)
	if err != nil {
		glog.Error(err)
		return errors.Wrap(err, "")
	}

	hs := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.imageCredentialInfo.SecretName,
			Namespace: ns,
		},
		Type: v1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			v1.DockerConfigJsonKey: bs,
		},
	}

	_, err = k.clientset.CoreV1().Secrets(ns).Create(ctx, hs, metav1.CreateOptions{})
	if err != nil {
		glog.Error(err)
		return errors.Wrap(err, "")
	}

	return nil
}

func (k *NamespaceWatcher) updateServiceAccount(ctx context.Context, ns string) error {
	if !k.isWatchNamespace(ns) {
		return nil
	}

	glog.Infof("updateServiceAccount Namespace: %s", ns)

	for _, saName := range k.imageCredentialInfo.ServiceAccounts {
		glog.Infof("updateServiceAccount Update Namespace: %s, ServiceAccount: %s", ns, saName)

		sa, err := k.clientset.CoreV1().ServiceAccounts(ns).Get(ctx, saName, metav1.GetOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				glog.Infof("updateServiceAccount, not found ServiceAccount %s, retry...", saName)
				// sleep and retry
				time.Sleep(time.Second * 10)
				sa, err = k.clientset.CoreV1().ServiceAccounts(ns).Get(ctx, saName, metav1.GetOptions{})
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						glog.Infof("updateServiceAccount, not found ServiceAccount %s, retry...", saName)
						// sleep and retry
						time.Sleep(time.Second * 20)
						sa, err = k.clientset.CoreV1().ServiceAccounts(ns).Get(ctx, saName, metav1.GetOptions{})
						if err != nil {
							glog.Error("updateServiceAccount failed", err)
							continue
						}
					} else {
						glog.Error(err)
						continue
					}
				}
			} else {
				glog.Error(err)
				continue
			}
		}

		had := false
		for _, ips := range sa.ImagePullSecrets {
			if ips.Name == k.imageCredentialInfo.SecretName {
				// 已有，则忽略
				had = true
				break
			}
		}
		if had {
			glog.Infof("updateServiceAccount Update Namespace: %s, ServiceAccount: %s, ServiceAccount's imagePullSecrets had %s, will ignore update",
				ns, saName, k.imageCredentialInfo.SecretName)
			continue
		}

		sa.ImagePullSecrets = append(sa.ImagePullSecrets, v1.LocalObjectReference{Name: k.imageCredentialInfo.SecretName})
		_, err = k.clientset.CoreV1().ServiceAccounts(ns).Update(ctx, sa, metav1.UpdateOptions{})
		if err != nil {
			glog.Error(err)
			continue
		}
	}

	return nil
}

func (k *NamespaceWatcher) isWatchNamespace(ns string) bool {
	for _, n := range k.imageCredentialInfo.WatchNamespaces {
		if n == "*" || n == ns {
			return true
		}
	}

	return false
}
