package system

import (
	"encoding/json"
	"errors"
	"evolution/backend/common/config"
	"evolution/backend/common/req"
	"evolution/backend/common/resp"
	"evolution/backend/system/models"
	"time"

	"github.com/mitchellh/mapstructure"
)

// TODO need super token

func UserByName(name string) (user models.User, err error) {
	condition := models.User{
		Name: name,
	}
	users, err := UserList(condition)
	if err != nil {
		return
	}
	if len(users) == 0 {
		err = errors.New("user not exist")
		return
	}
	if len(users) >= 2 {
		err = errors.New("user more than 2")
		return
	}
	user = users[0]
	return
}

func UserList(condition models.User) (users []models.User, err error) {
	conf := config.Get().System.System
	url := conf.Host + ":" + conf.Port + "/" + conf.Prefix + "/" + conf.Version + "/user/list"
	params, err := json.Marshal(condition)
	if err != nil {
		return
	}
	res, err := req.Exec(req.PostMethod, url, params, 300*time.Millisecond)
	var ret resp.Response
	err = json.Unmarshal(res, &ret)
	if err != nil {
		return
	}
	if ret.Code != 0 {
		err = errors.New("system user list by condition api error:" + ret.Desc)
		return
	}
	users = make([]models.User, 0)
	err = mapstructure.Decode(ret.Data, &users)
	return
}

func UserBAMap() (users map[string]string, err error) {
	conf := config.Get().System.System
	url := conf.Host + ":" + conf.Port + "/" + conf.Prefix + "/" + conf.Version + "/user/list"
	res, err := req.Exec(req.GetMethod, url, nil, 300*time.Millisecond)
	var ret resp.Response
	err = json.Unmarshal(res, &ret)
	if err != nil {
		return
	}
	if ret.Code != 0 {
		err = errors.New("system user list api error:" + ret.Desc)
		return
	}
	users = make(map[string]string)
	var usersStruct []models.User
	err = mapstructure.Decode(ret.Data, &usersStruct)
	if err != nil {
		return
	}
	for _, one := range usersStruct {
		users[one.Name] = one.Password
	}
	return
}
