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
	BuildPage(*xorm.Session)
	BuildSort(*xorm.Session)
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

type Page struct {
	Size    int    `xorm:"-" structs:"-" json:"Size,omitempty"`
	Current int    `xorm:"-" structs:"-" json:"Current,omitempty"`
	Order   string `xorm:"-" structs:"-" json:"Order,omitempty"`
}

type Sort struct {
	By      string `xorm:"-" structs:"-" json:"By,omitempty"`
	Reverse bool   `xorm:"-" structs:"-" json:"Reverse,omitempty"`
}

type Model struct {
	CreatedAt time.Time `xorm:"created comment('创建时间')" structs:"created_at,omitempty"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" structs:"updated_at,omitempty"`
	Page      *Page     `xorm:"-" structs:"-" json:"Page,omitempty"`
	Sort      *Sort     `xorm:"-" structs:"-" json:"Sort,omitempty"`
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

func (m *Model) BuildPage(session *xorm.Session) {
	if (m.Page != nil) && (m.Page.Size > 0) {
		start := (m.Page.Current - 1) * m.Page.Size
		session = session.Limit(m.Page.Size, start)
	}
}

func (m *Model) BuildSort(session *xorm.Session) {
	if (m.Sort != nil) && (len(m.Sort.By) != 0) {
		if m.Sort.Reverse {
			session.Desc(m.Sort.By)
		} else {
			session.Asc(m.Sort.By)
		}
	}
}
