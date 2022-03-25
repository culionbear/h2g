package h2g

type Func func()interface{}

type Config struct {
	Handler		map[string]Func
	Service		string
	Method		string
}