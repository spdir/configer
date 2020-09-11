package configer

import (
	"testing"
	"time"
)

type testConfig struct {
	Service *testServiceInfo `ini:"web" json:"web" yaml:"web" conf:"web"`
	Redis   *testRedisInfo   `ini:"redis" yaml:"redis"`
	Mysql   *testMysqlInfo   `ini:"mysql" yaml:"mysql"`
}
type testServiceInfo struct {
	Port         int    `default:"8091" env:"HTTP_PORT"`
	AppMode      string `default:"development"`
	EnablePprof  bool
	LogLevel     string `default:"info"`
	LogSaveDay   int    `default:"7"`
	LogSplitTime int    `default:"24"`
	LogOutType   string `default:"json"`
	LogOutPath   string `default:"console"`
	StartTime    string `default:"func:StartTime" func:"StartTime"`
}

type testRedisInfo struct {
	Host     string `panic:"redis host not is empty" env:"REDIS_HOST"`
	Port     int    `default:"6379" env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
}

type testMysqlInfo struct {
	Host        string `panic:"mysql host not is empty" env:"MYSQL_HOST"`
	Port        int    `default:"3306"  env:"MYSQL_PORT"`
	User        string `panic:"mysql user not is empty"  env:"MYSQL_USER"`
	Password    string `panic:"mysql password not is empty" env:"MYSQL_PASSWORD"`
	DB          string `panic:"mysql db name not is empty" env:"MYSQL_DB" conf:"db"`
	EnableDebug bool
}

func (*DefSetCall) StartTime(value string) string {
	_, err := time.Parse("2006/01/02", value)
	if err != nil {
		return time.Now().Format("2006/01/02")
	}
	return value
}

func (*DefGetCall) StartTime(value string) time.Time {
	runTime, err := time.Parse("2006/01/02", value)
	if err != nil {
		return time.Now()
	}
	return runTime
}

func TestLoadIniFile(t *testing.T) {
	testIniConfig := &testConfig{}
	err := LoadConfig(testIniConfig, "./test_config_file/test_config.ini")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", testIniConfig.Service)
}

func TestLoadJsonFile(t *testing.T) {
	testJSONConfig := &testConfig{}
	err := LoadConfig(testJSONConfig, "./test_config_file/test_config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", testJSONConfig.Mysql)
}

func TestLoadYamlFile(t *testing.T) {
	testYamlConfig := &testConfig{}
	err := LoadConfig(testYamlConfig, "./test_config_file/test_config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", testYamlConfig.Service)
}
