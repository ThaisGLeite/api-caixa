package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Logrus encapsula ferramentas de logging para diagnóstico e monitoramento.
type Logrus struct {
	// Logger é a instância do logger Logrus.
	Logger *logrus.Logger
}

// NewGoAppTools cria uma nova instância de GoAppTools com configurações padrão para o logger.
func NewGoAppTools() *Logrus {
	logger := logrus.New()
	logger.SetFormatter(new(fancyFormatter))
	return &Logrus{
		Logger: logger,
	}
}

func (app *Logrus) Info(msg string) {
	app.Logger.Info(msg)
}

// Check avalia se um erro ocorreu e, se verdadeiro, registra o erro usando Logrus.
// Diferente de outras implementações, essa função não encerra o programa, mas apenas registra o erro.
func (app *Logrus) Check(err error) {
	if err != nil {
		app.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("An error occurred:", err)
	}
}

// CheckAndPanic avalia se um erro ocorreu e, se verdadeiro, registra o erro e causa um panic.
func (app *Logrus) CheckAndPanic(err error) {
	if err != nil {
		app.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Panic("An error occurred:", err)
	}
}

type fancyFormatter struct{}

func (f *fancyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Customize your log format here
	return []byte(fmt.Sprintf("🌟 [Server Log] %s: %s\n", entry.Level, entry.Message)), nil
}
