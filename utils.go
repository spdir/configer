package configer

import (
	"errors"
	"reflect"
	"strconv"
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

func IsInSlice(data string, slice []string) bool {
	for _, s := range slice {
		if s == data {
			return true
		}
	}
	return false
}

func LowerCase(str string) string {
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
	confValue := value.(string)
	switch fieldType.Type.Kind() {
	case reflect.String:
		fieldVal.SetString(confValue)
	case reflect.Int:
		v, _ := strconv.Atoi(confValue)
		fieldVal.SetInt(int64(v))
	case reflect.Bool:
		v, _ := strconv.ParseBool(confValue)
		fieldVal.SetBool(v)
	default:
		return errors.New("undefiend set config type")
	}
	return
}
