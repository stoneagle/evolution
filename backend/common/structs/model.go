package structs

import (
	"reflect"
	"time"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

type JoinType string

func (jtype JoinType) String() string {
	return string(jtype)
}

const (
	InnerJoin JoinType = "INNER"
	LeftJoin  JoinType = "LEFT"
	RightJoin JoinType = "RIGHT"
)

type ModelGeneral interface {
	TableName() string
	Join() JoinGeneral
	BuildCondition(*xorm.Session)
	SlicePtr() interface{}
	Hook() error
	WithDeleted() bool
}

type JoinGeneral interface {
	Links() []JoinLinks
	SlicePtr() interface{}
	Transfer() ModelGeneral
	TransferCopy(ModelGeneral)
	TransferCopySlice(interface{}, interface{})
}

type JoinLinks struct {
	Table      string
	Type       JoinType
	LeftTable  string
	LeftField  string
	RightTable string
	RightField string
}

type Model struct {
	CreatedAt time.Time `xorm:"created comment('创建时间')" structs:"created_at,omitempty"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" structs:"updated_at,omitempty"`
}

type ModelWithId struct {
	Id    int   `xorm:"pk autoincr" structs:"id,omitempty"`
	Ids   []int `xorm:"-" structs:"id,omitempty" json:"Ids,omitempty"`
	Model `xorm:"extends"`
}

type ModelWithDeleted struct {
	ModelWithId `xorm:"extends"`
	DeletedAt   time.Time `xorm:"deleted comment('软删除时间')" structs:"deleted_at,omitempty"`
}

func (m *Model) BuildCondition(params map[string]interface{}, keyPrefix string) (condition builder.Eq) {
	condition = builder.Eq{}
	for key, value := range params {
		keyType := reflect.ValueOf(value).Kind()
		if keyType == reflect.Map {
			if key == "ModelWithDeleted" {
				valueMap := value.(map[string]interface{})
				// deletedAt, ok := valueMap["deleted_at"]
				// emptyTime := time.Time{}
				// if ok {
				// 	if deletedAt != emptyTime {
				// 		condition[keyPrefix+"deletedAt"] = id
				// 	}
				// }
				idMap, ok := valueMap["ModelWithId"]
				if ok {
					id, ok := idMap.(map[string]interface{})["id"]
					if ok {
						condition[keyPrefix+"id"] = id
					}
				}
			}
			continue
		}
		if keyType == reflect.Struct {
			continue
		}
		if keyType == reflect.String && value == "" {
			continue
		}
		if keyType == reflect.Int && value == 0 {
			continue
		}
		condition[keyPrefix+key] = value
	}
	return condition
}

func (m *Model) Join() JoinGeneral {
	return nil
}

func (m *Model) Hook() error {
	return nil
}

func (m *Model) WithDeleted() bool {
	return false
}
