package database

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Pedrommb91/go-auth/pkg/reflection"
)

type Tag string

func (t Tag) String() string {
	return string(t)
}

const (
	name      Tag = "name"
	reference Tag = "reference"
)

type modelParser[T any] struct {
	t any
}

func NewModelParser[T any](t any) modelParser[T] {
	return modelParser[T]{
		t: t,
	}
}

func (p modelParser[T]) GetTableName() string {
	// if its a pointer it erases the pointer signal
	return strings.ToLower(strings.Replace(reflection.GetType(p.t), "*", "", 1))
}

func (p modelParser[T]) GetColumns() []string {
	return reflection.GetAllTagsWithName(p.t, name.String())
}

func (p modelParser[T]) GetQueryColumns() string {
	return strings.Join(p.GetColumns()[:], ", ")
}

func (p modelParser[T]) GetValues() []string {
	values := reflection.GetAllValues(p.t)
	params := make([]string, 0)

	for _, v := range values {
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.String:
			if v == "" {
				params = append(params, "default")
			} else {
				params = append(params, fmt.Sprintf("'%s'", v.(string)))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.(int64) == 0 {
				params = append(params, "default")
			} else {
				values = append(values, strconv.FormatInt(v.(int64), 10))
			}
		default:
			params = append(params, "default")
		}
	}
	return params
}

func (p modelParser[T]) GetQueryValues() string {
	return strings.Join(p.GetValues()[:], ", ")
}

func (p modelParser[T]) HasRelations() bool {
	return len(reflection.GetAllTagsWithName(p.t, reference.String())) > 0
}

func (p modelParser[T]) GetAllRelationalStructs() []any {
	return reflection.GetAllStructsWithTagName(p.t, reference.String())
}

func (p modelParser[T]) GetTagNameByTypeName(field string) string {
	return reflection.GetTagByTypeName(p.t, field, name.String())
}
