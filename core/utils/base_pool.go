package utils

type BasePool interface {
	Run(i interface{})
	Shutdown()
}
