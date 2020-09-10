package configer

// different confiure type interface
type configInterface interface {
	loadFromFile(i interface{}, file string) error
	loadFromByte(i interface{}, data []byte) error
	getvalueFromSource(key string) (value interface{}, err error)
}
