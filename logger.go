package appcore

type Logger interface {
	Errorf(err error, format string, args ...interface{})
	Infof(format string, args ...interface{})
}
