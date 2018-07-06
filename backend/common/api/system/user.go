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

func UserBAMap(conf config.System) (users map[string]string, err error) {
	url := conf.Host + "/" + conf.Prefix + "/" + conf.Version + "/user/list"
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
	mapstructure.Decode(ret.Data, &usersStruct)
	for _, one := range usersStruct {
		users[one.Name] = one.Password
	}
	return
}
