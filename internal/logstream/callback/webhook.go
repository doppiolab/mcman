package callback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/doppiolab/mcman/internal/config"
	"github.com/doppiolab/mcman/internal/logstream"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// WebhookCallback sends log messages to a webhook.
//
// This callback will debonce the sending action to avoid rate limiting
func NewWebhookCallback(cfg *config.LogWebhookConfig) func(*logstream.LogBlock) error {
	debouncer := debounce.New(time.Duration(cfg.DebounceThreshold) * time.Millisecond)

	callback := &webhookCallback{
		cfg:       cfg,
		debouncer: debouncer,
	}

	return callback.OnLog
}

type webhookCallback struct {
	cfg           *config.LogWebhookConfig
	logLineBuffer []*logstream.LogBlock
	debouncer     func(func())
}

func (c *webhookCallback) OnLog(logBlock *logstream.LogBlock) error {
	c.logLineBuffer = append(c.logLineBuffer, logBlock)
	c.debouncer(c.fireWebhook)

	return nil
}

func (c *webhookCallback) fireWebhook() {
	if len(c.logLineBuffer) == 0 {
		return
	}

	// NOTE(jeongukjae): replace default client if required
	// TODO(jeongukjae): change ctx to context.WithTimeout
	ctx := context.Background()

	if c.cfg.SlackURL != "" {
		log.Debug().Int("n-logs", len(c.logLineBuffer)).Msg("send slack webhook")
		err := executeSlackWebhook(ctx, http.DefaultClient, c.cfg.SlackURL, c.logLineBuffer)
		if err != nil {
			log.Error().Err(err).Msg("failed to send discord webhook")
		}
	} else if c.cfg.DiscordURL != "" {
		log.Debug().Int("n-logs", len(c.logLineBuffer)).Msg("send discord webhook")
		err := executeDiscordWebhook(ctx, http.DefaultClient, c.cfg.DiscordURL, c.logLineBuffer)
		if err != nil {
			log.Error().Err(err).Msg("failed to send discord webhook")
		}
	}

	// empty buffer
	c.logLineBuffer = nil
}

type discordWebhookPayload struct {
	Content string `json:"content"`
}

type slackWebhookPayload struct {
	Text string `json:"text"`
}

func executeDiscordWebhook(
	ctx context.Context,
	client *http.Client,
	url string,
	logs []*logstream.LogBlock) error {
	messages := make([]string, len(logs))
	for i, log := range logs {
		messages[i] = log.String()
	}

	payload := &discordWebhookPayload{
		Content: fmt.Sprintf("```\n%s\n```", strings.Join(messages, "\n")),
	}

	// NOTE(jeogukjae): discord webhook server returns 204 if success
	return executeWebhook(ctx, client, url, payload, 204)
}

func executeSlackWebhook(
	ctx context.Context,
	client *http.Client,
	url string,
	logs []*logstream.LogBlock) error {
	messages := make([]string, len(logs))
	for i, log := range logs {
		messages[i] = log.String()
	}

	payload := &slackWebhookPayload{
		Text: fmt.Sprintf("```\n%s\n```", strings.Join(messages, "\n")),
	}

	return executeWebhook(ctx, client, url, payload, 200)
}

func executeWebhook(
	ctx context.Context,
	client *http.Client,
	url string,
	payload interface{},
	successCode int) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "cannot marshal json data")
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "cannot send request")
	}

	defer response.Body.Close()

	if response.StatusCode != successCode {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Error().Err(err).Msg("cannot read request body")
			return fmt.Errorf("non-%d response and cannot read request body. status: %s", successCode, response.Status)
		}

		return fmt.Errorf("non-%d response. status: %s, body: %s", successCode, response.Status, string(body))
	}

	return nil
}
