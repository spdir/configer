package configer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/tidwall/gjson"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

// different confiure type interface
type configInterface interface {
	loadConfig(i interface{}, file string) error
	getValue(key string) (value interface{}, err error)
}

/* configer different configer type module */
func initialConfigCallObject(configType string) (err error) {
	switch configType {
	case "ini":
		configCallObject = &iniConfig{}
	case "yaml":
		configCallObject = &yamlConfig{}
	case "json":
		configCallObject = &jsonConfig{}
	case "toml":
		configCallObject = &tomlConfig{}
	default:
		err = fmt.Errorf("not found %s config type module\n", configType)
		return
	}
	return nil
}

// ini configure
type iniConfig struct {
	configer *ini.File
}

func (conf *iniConfig) loadConfig(i interface{}, file string) error {
	iniFile, err := ini.Load(file)
	if err != nil {
		return err
	}
	conf.configer = iniFile
	err = iniFile.MapTo(i)
	if err != nil {
		return err
	}
	return nil
}

func (conf *iniConfig) getValue(key string) (value interface{}, err error) {
	return
}

// json configure
type jsonConfig struct {
	configer gjson.Result
}

func (conf *jsonConfig) loadConfig(i interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}
	conf.configer = gjson.ParseBytes(data)
	return nil
}

func (conf *jsonConfig) getValue(key string) (value interface{}, err error) {
	return
}

// yaml configure
type yamlConfig struct {
	configer gjson.Result
}

func (conf *yamlConfig) loadConfig(i interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, i)
	if err != nil {
		return err
	}
	conf.configer = gjson.ParseBytes(data)
	return nil
}

func (conf *yamlConfig) getValue(key string) (value interface{}, err error) {
	return
}

// toml configure
type tomlConfig struct {
	configer gjson.Result
}

func (conf *tomlConfig) loadConfig(i interface{}, file string) error {
	_, err := toml.DecodeFile(file, i)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	conf.configer = gjson.ParseBytes(data)
	return nil
}

func (conf *tomlConfig) getValue(key string) (value interface{}, err error) {
	return
}

/* configer parser */
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
			_ = setConfigValue(fieldType, fieldVal, envValue)
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
				_ = setConfigValue(fieldType, fieldVal, callResultValue)
			}
		} else {
			_ = setConfigValue(fieldType, fieldVal, defaultValue)
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

/* configer getter */
func readConfigVal(fileConfig *ini.File, keys []string) interface{} {
	sectionStrings := fileConfig.SectionStrings()
	isHas := isInSlice(keys[0], sectionStrings)
	if !isHas {
		return nil
	}
	section := fileConfig.Section(keys[0])
	keyStrings := section.KeyStrings()
	isHas = isInSlice(keys[1], keyStrings)
	if !isHas {
		return nil
	}
	value := section.Key(keys[1])
	return value
}

func getNextFloorConfig(config interface{}, key string) interface{} {
	tp := reflect.TypeOf(config)
	val := reflect.ValueOf(config)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		val = val.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		fieldTp := tp.Field(i)
		fieldVal := val.Field(i)
		fieldConfName := fieldTp.Tag.Get("conf")
		if fieldConfName == "" {
			fieldConfName = lowerCase(fieldTp.Name)
		}
		if key == fieldConfName {
			funcValue := fieldTp.Tag.Get("func")
			if funcValue != "" {
				callBack := reflect.ValueOf(&DefGetCall{})
				callFunc := callBack.MethodByName(funcValue)
				funcParamsLen := callFunc.Type().NumIn()
				var args []reflect.Value
				if funcParamsLen == 1 {
					args = append(args, reflect.ValueOf(fieldVal.Interface()))
				}
				callValue := callFunc.Call(args)
				if len(callValue) > 0 {
					return callValue[0].Interface()
				}
			}
			return fieldVal.Interface()
		}
	}
	return nil
}

// get configure field value
// example: Get("web::port")
// example: Get("config")
func Get(key string) (value interface{}) {
	if key == "" {
		return nil
	}
	keys := strings.Split(key, "::")
	if len(keys) == 0 {
		return nil
	}
	var config interface{}
	config = configerData
	for _, currentKey := range keys {
		config = getNextFloorConfig(config, currentKey)
		if config == nil {
			break
		}
	}
	if config == nil {
		config, err := configCallObject.getValue(key)
		if err != nil {
			return nil
		} else {
			return config
		}
	}
	return config
}

func GetString(key string) (value string) {
	result := Get(key)
	if v, ok := result.(string); ok {
		return v
	}
	return
}

func GetInt(key string) (value int) {
	result := Get(key)
	if v, ok := result.(int); ok {
		return v
	}
	return
}

func GetFloat64(key string) (value float64) {
	result := Get(key)
	if v, ok := result.(float64); ok {
		return v
	}
	return
}

func GetBool(key string) (value bool) {
	result := Get(key)
	if v, ok := result.(bool); ok {
		return v
	}
	return
}

func GetIntSlice(key string) (value []int) {
	result := Get(key)
	if v, ok := result.([]int); ok {
		return v
	}
	return
}

func GetStringMap(key string) (value map[string]interface{}) {
	result := Get(key)
	if v, ok := result.(map[string]interface{}); ok {
		return v
	}
	return
}

func GetStringMapString(key string) (value map[string]string) {
	result := Get(key)
	if v, ok := result.(map[string]string); ok {
		return v
	}
	return
}

func GetStringSlice(key string) (value []string) {
	result := Get(key)
	if v, ok := result.([]string); ok {
		return v
	}
	return
}

// func GetTime(key string) (value time.Time) {
// 	return
// }

// func GetDuration(key string) (value time.Duration) {
// return
// }

func IsSet(key string) (status bool) {
	val := Get(key)
	if val != nil {
		return true
	}
	return false
}

func AllSettings() (result map[string]interface{}) {
	data, err := json.Marshal(configerData)
	if err != nil {
		return
	}
	result = make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}
	return
}

/* configer update value */
// example: Set("web::port", 80)
// example: Set("title", "Configer Title")
func Set(key string, val interface{}) (err error) {
	if key == "" {
		return errors.New("update configer key params is null")
	}
	var lastKey string
	var rangeKeys []string
	keys := strings.Split(key, "::")
	if len(keys) == 1 {
		lastKey = key
	} else {
		lastKey = keys[len(keys)-1]
		rangeKeys = keys[0 : len(keys)-1]
	}
	var config interface{}
	config = configerData
	if len(rangeKeys) > 0 {
		for _, currentKey := range rangeKeys {
			lastConfig := getNextFloorConfig(config, currentKey)
			if config == nil {
				break
			}
			if checkObjIsStruct(lastConfig) {
				config = lastConfig
			} else {
				return fmt.Errorf("not found %s key\n", currentKey)
			}
		}
	}
	lastObj := getNextFloorConfig(config, lastKey)
	if lastObj == nil || checkObjIsStruct(lastObj) {
		return fmt.Errorf("not found %s key\n", lastKey)
	}
	tpField, vlField := getStructField(config, lastKey)
	if tpField != nil || vlField != nil {
		err = setConfigValue(tpField, vlField, val)
		if err != nil {
			return
		}
	}
	return
}

func getStructField(config interface{}, key string) (fieldType *reflect.StructField, fieldVal *reflect.Value) {
	tp := reflect.TypeOf(config)
	val := reflect.ValueOf(config)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		val = val.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		fieldTp := tp.Field(i)
		fieldVal := val.Field(i)
		fieldConfName := fieldTp.Tag.Get("conf")
		if fieldConfName == "" {
			fieldConfName = lowerCase(fieldTp.Name)
		}
		if key == fieldConfName {
			return &fieldTp, &fieldVal
		}
	}
	return nil, nil
}
