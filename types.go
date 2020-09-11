package configer

// different confiure type interface
type configInterface interface {
	loadConfig(i interface{}, file string) error
	getValue(key string) (value interface{}, err error)
}
