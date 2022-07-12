package logstream

import (
	"fmt"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/rs/zerolog/log"
)

type LogStream interface {
	Start()
	Stop()

	RegisterLogCallback(id string, callback func(string) error)
	DeregisterLogCallback(id string)
	GetNumLogCallbacks() int
}

type logStream struct {
	cfg          *config.LogWebhookConfig
	chans        map[string]chan string
	logCallbacks map[string]func(string) error
	quit         chan bool
}

func New(cfg *config.LogWebhookConfig, chans map[string]chan string) LogStream {
	return &logStream{
		cfg:          cfg,
		chans:        chans,
		logCallbacks: map[string]func(string) error{},
		quit:         make(chan bool),
	}
}

func (l *logStream) Start() {
	go l.streamChansToCallbacks()
	// TODO(hayeon): send debounced logs to webhook
}

func (l *logStream) Stop() {
	l.quit <- true
	l.sendAllRemainedData()
}

func (l *logStream) RegisterLogCallback(id string, callback func(string) error) {
	if _, ok := l.logCallbacks[id]; ok {
		log.Warn().Str("id", id).Msg("Log callback already registered. It will override the previous one.")
	}

	l.logCallbacks[id] = callback
}

func (l *logStream) DeregisterLogCallback(id string) {
	delete(l.logCallbacks, id)
}

func (l *logStream) GetNumLogCallbacks() int {
	return len(l.logCallbacks)
}

func (l *logStream) streamChansToCallbacks() {
	for {
		// l.chans can contain zero elements,
		// so we need to check l.quit channel first separately.
		select {
		case <-l.quit:
			return
		default:
		}

		for chanId, ch := range l.chans {
			select {
			case msg, ok := <-ch:
				if ok {
					l.sendToAllCallbacks(chanId, msg)
				}
			default:
			}
		}
	}
}

func (l *logStream) sendAllRemainedData() {
	for chanId, ch := range l.chans {
		l.sendReaminedData(chanId, ch)
	}
}

func (l *logStream) sendReaminedData(chanId string, ch chan string) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			l.sendToAllCallbacks(chanId, msg)
		default:
			return
		}
	}
}

func (l *logStream) sendToAllCallbacks(chanId, msg string) {
	for _, callback := range l.logCallbacks {
		if err := callback(fmt.Sprintf("[%s] %s", chanId, msg)); err != nil {
			log.Error().Err(err).Msg("Failed to call log callback")
		}
	}
}
