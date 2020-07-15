package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/airbloc/flux-slack-alert/slack"
	"github.com/airbloc/flux-slack-alert/webhook"
	"github.com/airbloc/logger"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	log := logger.New("fluxslack")

	cfg := config{}
	if err := envconfig.Process("", &cfg); err != nil {
		log.Error("failed to parse config", err)
		os.Exit(1)
	}

	sender, err := slack.NewSender(&cfg.Config)
	if err != nil {
		log.Error("failed to initialize Slack sender", err)
		os.Exit(1)
	}

	w := webhook.New(cfg.Port, sender)
	if err := w.Start(); err != nil {
		log.Error("failed to start webhook server", err)
		os.Exit(1)
	}
	log.Info("started forwarding flux event to {}", cfg.SlackWebhookURL)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := w.Close(); err != nil {
		log.Error("error occurred while shutting down the webhook server", err)
		os.Exit(1)
	}
	log.Info("bye")
}
