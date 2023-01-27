package configuration

type Configuration interface {
	Update() *Configuration
	Get(name string) (value string)
	GetInt(name string) (value int64)
	GetBool(name string) (value bool)
	read() error
	write()
}
