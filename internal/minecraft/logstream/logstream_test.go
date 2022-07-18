package logstream

import (
	"testing"
	"time"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogStream_Callback(t *testing.T) {
	results := []string{}
	resultsAppendCallback := func(l *LogBlock) error {
		results = append(results, l.String())
		return nil
	}

	testChan := make(chan string)
	chans := map[string]chan string{"test": testChan}
	logStream := New(&config.LogWebhookConfig{}, chans)

	logStream.Start()
	defer logStream.Stop()

	// Register callback and send dummy text
	logStream.RegisterLogCallback("test", resultsAppendCallback)
	testChan <- "dummy_text"
	sleep5ms(t)

	require.Equal(t, 1, len(results))
	require.Equal(t, "[test] dummy_text", results[0])

	// Deregister callback and send dummy text
	logStream.DeregisterLogCallback("test")
	testChan <- "dummy_text"
	sleep5ms(t)

	require.Equal(t, 1, len(results)) // it should not be appended anymore
}

func TestLogStream_CallbackWithMultipleChannel(t *testing.T) {
	results := []string{}
	resultsAppendCallback := func(l *LogBlock) error {
		results = append(results, l.String())
		return nil
	}

	testChan := make(chan string)
	testChan2 := make(chan string)
	chans := map[string]chan string{"test": testChan, "test2": testChan2}
	logStream := New(&config.LogWebhookConfig{}, chans)

	logStream.Start()
	defer logStream.Stop()

	// Register callback and send dummy text
	logStream.RegisterLogCallback("test", resultsAppendCallback)
	testChan <- "dummy_text"
	testChan2 <- "dummy_text 123"
	sleep5ms(t)

	require.Equal(t, 2, len(results))
	require.Equal(t, "[test] dummy_text", results[0])
	require.Equal(t, "[test2] dummy_text 123", results[1])
}

func TestLogStream_TryToDeregisterUnknown(t *testing.T) {
	emptyCallback := func(*LogBlock) error { return nil }

	chans := map[string]chan string{}
	logStream := New(&config.LogWebhookConfig{}, chans)

	logStream.Start()
	defer logStream.Stop()

	// Try to deregister unregistered callback id
	logStream.RegisterLogCallback("test", emptyCallback)
	logStream.DeregisterLogCallback("test2")

	assert.Equal(t, 1, logStream.GetNumLogCallbacks())
}

func TestLogStream_SendAfterStop(t *testing.T) {
	emptyCallback := func(*LogBlock) error { return nil }

	testChan := make(chan string)
	chans := map[string]chan string{"test": testChan}
	logStream := New(&config.LogWebhookConfig{}, chans)

	logStream.Start()
	logStream.RegisterLogCallback("test", emptyCallback)
	logStream.Stop()

	// Try to send log after stopping logStream
	select {
	case testChan <- "dummy_text":
		t.Error("Should not be able to send message after stop")
	default:
		t.Log("success")
	}
}

func sleep5ms(t *testing.T) {
	t.Helper()
	time.Sleep(time.Millisecond * 5)
}
