package configer

import "errors"

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
