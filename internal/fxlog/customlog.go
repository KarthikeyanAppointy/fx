package fxlog

import (
	"fmt"
	"go.uber.org/fx/internal/fxreflect"
	"go.uber.org/zap"
	"os"
	"strings"
)

type CustomLogger struct {
	Logger *zap.Logger
}

type LogRequest struct {
	Message string
	Objects []interface{}
}

func NewCustomLogger(logger *zap.Logger) CustomLogger {
	return CustomLogger{
		Logger: logger,
	}
}

func (l CustomLogger) Printf(format string, v ...interface{}) {
	l.Log(&LogRequest{
		Message: format,
		Objects: v,
	})
}

// PrintProvide logs a type provided into the dig.Container.
func (l CustomLogger) PrintProvide(t interface{}) {
	for _, rtype := range fxreflect.ReturnTypes(t) {
		l.Printf("PROVIDE\t%s <= %s", rtype, fxreflect.FuncName(t))
	}
}

func (l CustomLogger) PrintDecorate(t interface{}) {
	for _, rtype := range fxreflect.ReturnTypes(t) {
		l.Printf("DECORATE\t%s <= %s", rtype, fxreflect.FuncName(t))
	}
}

// PrintSignal logs an os.Signal.
func (l CustomLogger) PrintSignal(signal os.Signal) {
	l.Printf(strings.ToUpper(signal.String()))
}

func (l CustomLogger) Panic(err error) {
	l.Printf(prepend(err.Error()))
	panic(err)
}

func (l CustomLogger) Log(in *LogRequest) {

	l.Logger.Info(fmt.Sprintf(in.Message, in.Objects...))

}

// Fatalf logs an Fx line then fatals.
func (l CustomLogger) Fatalf(format string, v ...interface{}) {
	l.Printf(prepend(format), v...)
	_exit()
}
