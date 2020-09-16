package configer

import (
	"errors"
	"reflect"
	"strings"
)

func getFileType(file string) string {
	if strings.HasSuffix(file, ".ini") {
		return "ini"
	}
	if strings.HasSuffix(file, ".toml") {
		return "toml"
	}
	if strings.HasSuffix(file, ".yaml") || strings.HasPrefix(file, ".yml") {
		return "yaml"
	}
	if strings.HasSuffix(file, ".json") {
		return "json"
	}
	return ""
}

func isInSlice(data string, slice []string) bool {
	for _, s := range slice {
		if s == data {
			return true
		}
	}
	return false
}

func lowerCase(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 {
				vv[i] += 32
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func checkObjIsStruct(i interface{}) bool {
	tp := reflect.TypeOf(i)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	if tp.Kind() == reflect.Struct {
		return true
	}
	return false
}

func setConfigValue(fieldType *reflect.StructField, fieldVal *reflect.Value, value interface{}) (err error) {
	err = errors.New("undefiend type")
	switch fieldType.Type.Kind() {
	case reflect.String:
		if v, ok := value.(string); ok {
			fieldVal.SetString(v)
		} else {
			return
		}
	case reflect.Int, reflect.Int64:
		if v, ok := value.(int); ok {
			fieldVal.SetInt(int64(v))
		} else {
			return
		}
	case reflect.Bool:
		if v, ok := value.(bool); ok {
			fieldVal.SetBool(v)
		} else {
			return
		}
	case reflect.Float64:
		fieldVal.SetFloat(value.(float64))
	default:
		return errors.New("undefiend set config type")
	}
	return nil
}
