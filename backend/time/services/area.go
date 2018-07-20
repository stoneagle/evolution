package services

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"

	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"fmt"
	"strconv"
)

type Area struct {
	Pack ServicePackage
	structs.Service
}

func NewArea(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Area {
	ret := Area{}
	ret.Init(engine, cache, log)
	return &ret
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
