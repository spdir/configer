package configer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

func initialConfigCallObject(configType string) (err error) {
	fmt.Println("type", configType)
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
		err = fmt.Errorf("not found %s cofnig type module\n", configType)
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
}

func (conf *tomlConfig) loadConfig(i interface{}, file string) error {

	return nil
}

func (conf *tomlConfig) getValue(key string) (value interface{}, err error) {
	return
}
