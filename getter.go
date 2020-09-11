package configer

import "time"

// get configure field value
// example: Get("web::port")
// example: Get("config")
func Get(key string) (value interface{}) {
	return nil
}

func GetString(key string) (value string) {
	return
}

func GetInt(key string) (value int) {
	return
}

func GetFloat64(key string) (value float64) {
	return
}

func GetBool(key string) (value bool) {
	return
}

func GetIntSlice(key string) (value []int) {
	return
}

func GetStringMap(key string) (value map[string]interface{}) {
	return
}

func GetStringMapString(key string) (value map[string]string) {
	return
}

func GetStringSlice(key string) (value []string) {
	return
}

func GetTime(key string) (value time.Time) {
	return
}

func GetDuration(key string) (value time.Duration) {
	return
}

func IsSet(key string) (status bool) {
	return
}

func AllSettings() (result map[string]interface{}) {
	return
}
