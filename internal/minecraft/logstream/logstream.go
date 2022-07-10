package logstream

import "github.com/doppiolab/mcman/internal/config"

type LogStream interface {
	Start() error
	Stop() error

	RegisterLogCallback(id string, callback func(string) error) error
	DeregisterLogCallback(id string)
}

type logStream struct {
	cfg   *config.LogWebhookConfig
	chans map[string]<-chan string

	logCallbacks map[string]func(string) error
}

func NewLogstream(cfg *config.LogWebhookConfig, chans map[string]<-chan string) (LogStream, error) {
	return &logStream{cfg: cfg, chans: chans}, nil
}

func (l *logStream) Start() error {}

func (l *logStream) Stop() error {}

func (l *logStream) RegisterLogCallback(id string, callback func(string) error) error {}

func (l *logStream) DeregisterLogCallback(id string) {}
