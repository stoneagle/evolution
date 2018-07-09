package services

import (
	"errors"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Area struct {
	structs.Service
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

func (s *Area) ListWithCondition(area *models.Area) (areas []models.Area, err error) {
	areas = make([]models.Area, 0)
	sql := s.Engine.Asc("parent")
	condition := area.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&areas)
	if err != nil {
		return
	}
	return
}

func (s *Area) List() (areas []models.Area, err error) {
	areas = make([]models.Area, 0)
	err = s.Engine.Asc("parent").Find(&areas)
	return
}

func (s *Area) TransferListToTree(areas []models.Area, fieldMap map[int]string) (areaTrees map[int]models.AreaTree, err error) {
	areaTrees = make(map[int]models.AreaTree)

buildTreeLoop:
	leftArea := make([]models.Area, 0)
	for _, one := range areas {
		node := models.AreaNode{
			Id:       one.Id,
			Value:    one.Name,
			Children: make([]models.AreaNode, 0),
		}

		fieldName, exist := fieldMap[one.FieldId]
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
			var targetFlag bool
			tree.Children, targetFlag = locationNode(tree.Children, node, one)
			if !targetFlag {
				leftArea = append(leftArea, one)
			}
		}
		areaTrees[one.FieldId] = tree
	}

	if len(leftArea) > 0 {
		areas = leftArea
		goto buildTreeLoop
	}
	return
}

func (s *Area) GetAllLeafId(areaId int) (leafIds []int, err error) {
	// this version can not solve parentId > id situation
	// subSql := fmt.Sprintf("SELECT id FROM ( "+
	// 	"	SELECT  id,"+
	// 	"					NAME,"+
	// 	"					parent,"+
	// 	"					TYPE"+
	// 	"	FROM    (SELECT * FROM `area`"+
	// 	"					 ORDER BY parent, id) products_sorted,"+
	// 	"					(SELECT @pv := %v) initialisation"+
	// 	"	WHERE   FIND_IN_SET(parent, @pv)"+
	// 	"	AND     LENGTH(@pv := CONCAT(@pv, ',', id))"+
	// 	") result WHERE TYPE = %v ORDER BY parent;", areaId, models.AreaTypeLeaf)
	subSql := fmt.Sprintf("WITH recursive cte (id, NAME, parent, TYPE) AS"+
		"("+
		" SELECT     id,"+
		"            NAME,"+
		"            parent,"+
		"            TYPE"+
		" FROM       `area`"+
		" WHERE      parent = %v"+
		" UNION ALL"+
		" SELECT     p.id,"+
		"            p.name,"+
		"            p.parent,"+
		"            p.type"+
		" FROM       `area` p"+
		" INNER JOIN cte"+
		"         ON p.parent = cte.id"+
		")"+
		"SELECT * FROM cte WHERE TYPE = %v ORDER BY parent;", areaId, models.AreaTypeLeaf)
	results, err := s.Engine.Query(subSql)
	if err != nil {
		return leafIds, err
	}
	leafIds = make([]int, 0)
	for _, one := range results {
		areaId, err := strconv.Atoi(string(one["id"]))
		if err != nil {
			return leafIds, err
		}
		leafIds = append(leafIds, areaId)
	}
	return
}

func locationNode(level []models.AreaNode, node models.AreaNode, one models.Area) ([]models.AreaNode, bool) {
	targetFlag := false
	for key, leaf := range level {
		if leaf.Id == one.Parent {
			level[key].Children = append(leaf.Children, node)
			targetFlag = true
			break
		} else if len(leaf.Children) > 0 {
			level[key].Children, targetFlag = locationNode(leaf.Children, node, one)
			if targetFlag {
				break
			}
		}
	}
	return level, targetFlag
}
