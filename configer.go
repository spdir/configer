package configer

import "errors"

// Please follow the format of each configuration file and use tag when struct field
/* struct tag describe
	config type tag: (ini, json, yaml, toml) 不同配置原始tag，用于struct字段和配置文件不同时读取
	conf: 获取值时定义的tag，默认：以struct字段名首页字母小写
	default: 设置默认值, 如果字段时为空时（单字符串表示直接设置值，func:xxx表示调用函数设置值）
	env: 读取环境变量中的值，默认：env名为struct 字段名以大写字母分割，每个分割的单词全大写 以 _ 相连接
	func: 获取值时调用的函数声明
	none: 设置字段什么值为none，默认根据字段类空值去判断
	required: bool值，是否为必输字段
	panic: 当required 为 true 时 这个值如果为空, panic提示的错误消息内容
	pass: 忽略初始化扫描
---
	config type tag: (ini, json, yaml, toml) Different configuration original tags, used for reading when struct fields and configuration files are different
	conf: The tag defined when getting the value, default: lowercase the first letter of the struct field name
	default: Set the default value, if the field is empty (single string means directly set the value, func:xxx means call the function to set the value)
	env: Read the value in the environment variable, default: env is named struct, the field name is separated by uppercase letters, and each divided word is all uppercase and connected with _
	func: Function declaration called when getting the value
	none: Set the value of the field to none, by default it is judged based on the null value of the field class
	required: bool value, whether it is a required field
	panic: When required is true, if this value is empty, the content of the error message that panic prompts
	pass: ignore initial scan
---
*/

/* define data */
var supportConfig = []string{"ini", "yaml", "json", "toml"}

// different config bind interface
var configCallObject configInterface

// config
var configer *config

type config struct {
	configFileType string
	configFile     string
}

// binding global configure object
var configerData interface{}

// default parse set value callback funcation binding object
// The function has one and only one value parameter, but no parameter, the first letter of the function name should be capitalized
// example: func (*DefSetCall) StartTime(value string) string {...}
type DefSetCall struct{}

// default get configure value callback fucnation binding object
// The function has one and only one value parameter, but no parameter, the first letter of the function name should be capitalized
// example: func (*DefGetCall) StartTime(value string) time.Time {...}
type DefGetCall struct{}

/* configer main */
func LoadConfig(i interface{}, file string) (err error) {
	err = new(i, file)
	if err != nil {
		return
	}
	return nil
}

func new(i interface{}, file string) (err error) {
	fileType := getFileType(file)
	if !IsInSlice(fileType, supportConfig) {
		err = errors.New("not support config type")
		return
	}

	configer = &config{
		configFile: file,
	}
	configer.configFileType = fileType
	configerData = i

	err = initialConfigCallObject(configer.configFileType)
	if err != nil {
		return
	}

	// load different configure type
	err = configCallObject.loadConfig(configerData, configer.configFile)
	if err != nil {
		return
	}

	// parser different configure type
	parserConfig(configerData)

	return nil
}
