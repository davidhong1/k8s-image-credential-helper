package service

import (
	"encoding/base64"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type DockerConfigJson struct {
	Auths map[string]DockerRegistry `json:"auths"`
}

type DockerRegistry struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}

func dockerConfigJsonKeyBytes(Host, User, Password, email string) ([]byte, error) {
	dr := DockerRegistry{
		Username: User,
		Password: Password,
		Email:    email,
		Auth:     base64.StdEncoding.EncodeToString([]byte(User + ":" + Password)),
	}
	dcj := DockerConfigJson{
		Auths: map[string]DockerRegistry{Host: dr},
	}
	bs, err := json.Marshal(&dcj)
	if err != nil {
		glog.Error(err)
		return nil, errors.Wrap(err, "")
	}

	return bs, nil
}
