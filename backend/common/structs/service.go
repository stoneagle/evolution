package structs

import (
	"errors"
	"evolution/backend/common/logger"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type ServiceGeneral interface {
	One(int, interface{}) error
	List(interface{}) error
	Create(interface{}) error
	Update(int, interface{}) error
	Delete(int, interface{}) error
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

func (s *Service) One(id int, modelPtr interface{}) (err error) {
	session := s.Engine.Where("id = ?", id)
	has, err := session.Get(modelPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	if !has {
		return errors.New(fmt.Sprintf("%v not exist", s.Logger.Resource))
	}
	return
}

func (s *Service) Delete(id int, modelPtr interface{}) (err error) {
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

func (s *Service) List(modelsPtr interface{}) (err error) {
	session := s.Engine.Where("1 = 1")
	err = session.Limit(10).Find(modelsPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	return
}

func (s *Service) Create(modelPtr interface{}) (err error) {
	session := s.Engine.Where("1 = 1")
	_, err = session.Insert(modelPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	return
}

func (s *Service) Update(id int, modelPtr interface{}) (err error) {
	session := s.Engine.Id(id)
	_, err = session.Update(modelPtr)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	return
}
