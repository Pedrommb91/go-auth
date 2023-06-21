package reflection

import (
	"reflect"
)

func GetType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func GetAllFields(s interface{}) []string {
	fields := make([]string, 0)

	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}

	return fields
}

func GetAllTags(s interface{}) []string {
	tags := make([]string, 0)

	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		tags = append(tags, string(t.Field(i).Tag))
	}

	return tags
}

func GetAllTagsWithName(s interface{}, name string) []string {
	tags := make([]string, 0)

	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		tag := string(t.Field(i).Tag.Get(name))
		tags = append(tags, tag)
	}

	return tags
}

func GetAllStructsWithTagName(s interface{}, name string) []any {
	r := make([]any, 0)
	val := reflect.Indirect(reflect.ValueOf(s))

	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		tag := string(t.Field(i).Tag.Get(name))
		if tag != "" {
			valueField := val.Field(i)
			r = append(r, valueField.Interface())
		}
	}

	return r
}

func GetTagByTypeName(s interface{}, field, tag string) string {
	val := reflect.Indirect(reflect.ValueOf(s))
	for i := 0; i < val.NumField(); i++ {
		typeName := val.Field(i).Type().Name()
		if typeName == field {
			return val.Type().Field(i).Tag.Get(tag)
		}
	}
	return ""
}

func GetAllValuesAsString(s interface{}) []any {
	values := make([]any, 0)
	val := reflect.Indirect(reflect.ValueOf(s))

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)

		f := valueField.Interface()
		val := reflect.ValueOf(f)

		switch val.Kind() {
		case reflect.String:
			if val.Kind() == reflect.String {
				values = append(values, val.String())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values = append(values, val.Int())
		default:
			values = append(values, "")
		}

	}

	return values
}
