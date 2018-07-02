package services

import (
	"errors"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Area struct {
	Basic
}

func NewArea(engine *xorm.Engine, cache *redis.Client) *Area {
	ret := Area{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Area) One(id int) (interface{}, error) {
	model := models.Area{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Area) Add(model models.Area) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Area) UpdateByMap(id int, model map[string]interface{}) (err error) {
	_, err = s.Engine.Table(new(models.Area)).Id(id).Update(&model)
	return
}

func (s *Area) Update(id int, model models.Area) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Area) Delete(id int, model models.Area) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Area) List() (countries []models.Area, err error) {
	countries = make([]models.Area, 0)
	err = s.Engine.Where("del = ?", 0).Asc("parent").Find(&countries)
	return
}

func (s *Area) TransferListToTree(areas []models.Area) (areaTrees map[int]models.AreaTree, err error) {
	areaTrees = make(map[int]models.AreaTree)
	for _, one := range areas {
		node := models.AreaNode{
			Id:       one.Id,
			Value:    one.Name,
			Children: make([]models.AreaNode, 0),
		}

		fieldName, exist := models.AreaFiledMap[one.FieldId]
		if !exist {
			err = errors.New("field Id not exist")
			return
		}

		_, ok := areaTrees[one.FieldId]
		if !ok {
			areaTrees[one.FieldId] = models.AreaTree{
				Value:    fieldName,
				Children: make([]models.AreaNode, 0),
			}
		}

		tree := areaTrees[one.FieldId]
		if one.Parent == 0 {
			tree.Children = append(tree.Children, node)
		} else {
			tree.Children = locationNode(tree.Children, node, one)
		}
		areaTrees[one.FieldId] = tree
	}
	return
}

func locationNode(level []models.AreaNode, node models.AreaNode, one models.Area) []models.AreaNode {
	for key, leaf := range level {
		if leaf.Id == one.Parent {
			level[key].Children = append(leaf.Children, node)
			break
		} else if len(leaf.Children) > 0 {
			level[key].Children = locationNode(leaf.Children, node, one)
		}
	}
	return level
}
