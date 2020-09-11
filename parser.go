package configer

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// parser configure main
func parserConfig(config interface{}) {
	scanConfig(config)
}

func scanConfig(config interface{}) {
	tp := reflect.TypeOf(config)
	val := reflect.ValueOf(config)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		val = val.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		fieldTp := tp.Field(i)
		fieldVal := val.Field(i)
		var tyKind reflect.Kind
		passTag := fieldTp.Tag.Get("pass")
		if passTag != "" {
			continue
		}
		if fieldTp.Type.Kind() == reflect.Ptr {
			tyKind = fieldTp.Type.Elem().Kind()
		}
		if tyKind == reflect.Struct {
			scanConfig(fieldVal.Interface())
		} else {
			updateConfigValue(&fieldTp, &fieldVal)
		}
	}
}

func updateConfigValue(fieldType *reflect.StructField, fieldVal *reflect.Value) {
	name := fieldType.Name
	value := fieldVal.Interface()
	// env name
	envName := fieldType.Tag.Get("env")
	if envName == "" {
		reCompile := regexp.MustCompile("[A-Z]+[a-z0-9]+")
		keySlice := reCompile.FindAllString(name, -1)
		for i := 0; i < len(keySlice); i++ {
			keySlice[i] = strings.ToUpper(keySlice[i])
		}
		envName = strings.Join(keySlice, "_")
	}
	// default value
	defaultType := "" // "","func"
	defaultValue := fieldType.Tag.Get("default")
	if strings.HasPrefix(defaultValue, "func:") {
		defaultType = "func"
		defaultValue = strings.Split(defaultValue, "func:")[1]
	}
	// none value
	noneValue := fieldType.Tag.Get("none")
	if noneValue == "" {
		switch fieldType.Type.Kind() {
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			noneValue = "0"
		case reflect.Bool:
			noneValue = "false"
		}
	}

	// priority: env, source, default, none
	envValue := os.Getenv(envName)
	if envValue != "-" {
		if envValue != "" {
			setConfigValue(fieldType, fieldVal, envValue)
			return
		}
	}
	newVal := fmt.Sprintf("%v", value)
	if newVal != strings.TrimSpace(noneValue) {
		return
	}

	if defaultValue != "" {
		if defaultType == "func" {
			callBack := reflect.ValueOf(&DefSetCall{})
			callFunc := callBack.MethodByName(defaultValue)
			funcParamsLen := callFunc.Type().NumIn()
			var args []reflect.Value
			if funcParamsLen == 1 {
				args = append(args, reflect.ValueOf(newVal))
			}
			callValue := callFunc.Call(args)
			if len(callValue) > 0 {
				callResultValue := callValue[0].String()
				setConfigValue(fieldType, fieldVal, callResultValue)
			}
		} else {
			setConfigValue(fieldType, fieldVal, defaultValue)
		}
	}

	// error panic
	isRequired := fieldType.Tag.Get("required")
	if isRequired == "true" {
		if fieldVal.String() == "" && defaultValue == "" && envValue == "" {
			panicValue := fieldType.Tag.Get("panic")
			if panicValue != "" {
				panic(panicValue)
			} else {
				panic(fmt.Sprintf("%s is required field!!!", name))
			}
		}
	}
}

func setConfigValue(fieldType *reflect.StructField, fieldVal *reflect.Value, value interface{}) {
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
	}
}
