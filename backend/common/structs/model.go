package structs

import (
	"reflect"
	"time"

	"github.com/go-xorm/builder"
)

type ModelGeneral interface {
	BuildCondition() builder.Eq
}

type Model struct {
	CreatedAt time.Time `xorm:"created comment('创建时间')" structs:"created_at,omitempty"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" structs:"updated_at,omitempty"`
}

type ModelWithId struct {
	Id    int   `xorm:"pk autoincr" structs:"id,omitempty"`
	Ids   []int `xorm:"-" structs:"id,omitempty"`
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
		if keyType == reflect.Map || keyType == reflect.Struct {
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
