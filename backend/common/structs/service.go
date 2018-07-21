package structs

import (
	"errors"
	"evolution/backend/common/logger"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type ServiceGeneral interface {
	One(int, ModelGeneral) error
	List(ModelGeneral, interface{}) error
	Create(ModelGeneral) error
	Update(int, ModelGeneral) error
	UpdateByMap(string, int, map[string]interface{}) error
	Delete(int, ModelGeneral) error
}

type Service struct {
	Engine *xorm.Engine
	Cache  *redis.Client
	Logger *logger.Logger
}

func (s *Service) Init(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) {
	s.Engine = engine
	s.Cache = cache
	s.Logger = log
}

func (s *Service) SetEngine(engine *xorm.Engine) {
	s.Engine = engine
}

func (s *Service) SetCache(cache *redis.Client) {
	s.Cache = cache
}

func (s *Service) LogSql(sql string, args interface{}, err error) {
	info := fmt.Sprintf("[SQL] %v %#v", sql, args)
	s.Logger.Log(logger.WarnLevel, info, err)
}

func (s *Service) One(id int, modelPtr ModelGeneral) (err error) {
	joinPtr := modelPtr.Join()
	var has bool
	session := s.Engine.Unscoped().Table(modelPtr.TableName()).Where(fmt.Sprintf("%s.id = ?", modelPtr.TableName()), id)
	if joinPtr != nil {
		links := joinPtr.Links()
		for _, link := range links {
			session.Join(link.Type.String(), link.Table, fmt.Sprintf("%s.%s = %s.%s", link.LeftTable, link.LeftField, link.RightTable, link.RightField))
		}
		has, err = session.Get(joinPtr)
	} else {
		has, err = session.Get(modelPtr)
	}

	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	if !has {
		return errors.New(fmt.Sprintf("%v not exist", s.Logger.Resource))
	}

	if joinPtr != nil {
		joinPtr.TransferCopy(modelPtr)
	}
	return
}

func (s *Service) List(modelPtr ModelGeneral, modelsPtr interface{}) (err error) {
	join := modelPtr.Join()
	if join != nil {
		joinsPtr := join.SlicePtr()
		session := s.Engine.Unscoped().Table(modelPtr.TableName())
		links := join.Links()
		for _, link := range links {
			session.Join(link.Type.String(), link.Table, fmt.Sprintf("%s.%s = %s.%s", link.LeftTable, link.LeftField, link.RightTable, link.RightField))
		}
		modelPtr.BuildCondition(session)
		err = session.Find(joinsPtr)
		if err != nil {
			sql, args := session.LastSQL()
			s.LogSql(sql, args, err)
			return
		}
		join.TransferCopySlice(joinsPtr, modelsPtr)
	} else {
		session := s.Engine.Unscoped().Table(modelPtr.TableName())
		modelPtr.BuildCondition(session)
		err = session.Limit(10).Find(modelsPtr)
		if err != nil {
			sql, args := session.LastSQL()
			s.LogSql(sql, args, err)
		}
	}
	return
}

func (s *Service) Create(modelPtr ModelGeneral) (err error) {
	session := s.Engine.Where("1 = 1")
	_, err = session.Insert(modelPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	return
}

func (s *Service) Update(id int, modelPtr ModelGeneral) (err error) {
	session := s.Engine.Id(id)
	_, err = session.Update(modelPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	return
}

func (s *Service) UpdateByMap(table string, id int, params map[string]interface{}) (err error) {
	_, err = s.Engine.Table(table).Id(id).Update(&params)
	return
}

func (s *Service) Delete(id int, modelPtr ModelGeneral) (err error) {
	session := s.Engine.Id(id)
	has, err := session.Get(modelPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	if !has {
		return errors.New(fmt.Sprintf("%v not exist", s.Logger.Resource))
	}
	_, err = s.Engine.Id(id).Delete(modelPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	return
}
