package main

import "github.com/airbloc/flux-slack-alert/slack"

type config struct {
	slack.Config

	Port int `default:"8080"`
}
