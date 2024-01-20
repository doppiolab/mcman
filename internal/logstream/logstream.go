package logstream

import (
	"github.com/doppiolab/mcman/internal/config"
	"github.com/rs/zerolog/log"
)

type LogStream interface {
	Start()
	Stop()

	RegisterLogCallback(id string, callback func(*LogBlock) error)
	DeregisterLogCallback(id string)
	GetNumLogCallbacks() int
}

type logStream struct {
	cfg          *config.LogWebhookConfig
	chans        map[string]chan string
	logCallbacks map[string]func(*LogBlock) error
	quit         chan bool
}

func New(cfg *config.LogWebhookConfig, chans map[string]chan string) LogStream {
	return &logStream{
		cfg:          cfg,
		chans:        chans,
		logCallbacks: map[string]func(*LogBlock) error{},
		quit:         make(chan bool),
	}
}

func (l *logStream) Start() {
	for chanID, ch := range l.chans {
		go l.streamChansToCallbacks(chanID, ch)
	}
}

func (l *logStream) Stop() {
	for _, ch := range l.chans {
		_ = ch // Ignore the unused variable
		l.quit <- true
	}
	l.sendAllRemainedData()
}

func (l *logStream) RegisterLogCallback(id string, callback func(*LogBlock) error) {
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

func (l *logStream) streamChansToCallbacks(chanID string, ch chan string) {
	for {
		select {
		case <-l.quit:
			return
		case msg, ok := <-ch:
			if !ok {
				continue
			}

			logBlock := &LogBlock{
				ChanID: chanID,
				Msg:    msg,
			}

			// send message to callbacks
			l.sendToAllCallbacks(logBlock)
		}
	}
}

func (l *logStream) sendAllRemainedData() {
	for chanID, ch := range l.chans {
		l.sendRemainedData(chanID, ch)
	}
}

func (l *logStream) sendRemainedData(chanID string, ch chan string) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			logBlock := &LogBlock{
				ChanID: chanID,
				Msg:    msg,
			}
			l.sendToAllCallbacks(logBlock)
		default:
			return
		}
	}
}

func (l *logStream) sendToAllCallbacks(logBlock *LogBlock) {
	for _, callback := range l.logCallbacks {
		if err := callback(logBlock); err != nil {
			log.Error().Err(err).Msg("Failed to call log callback")
		}
	}
}
