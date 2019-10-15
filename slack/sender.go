package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/airbloc/logger"
	fluxevent "github.com/fluxcd/flux/pkg/event"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Sender interface {
	Compose(event fluxevent.Event) slack.Message
	Send(msg slack.Message) error
}

type sender struct {
	log    *logger.Logger
	config *Config

	resourceURITmpl *template.Template
}

func NewSender(config *Config) (Sender, error) {
	tmpl := template.New("resource-uri")
	tmpl, err := tmpl.Parse(config.WorkloadURITemplate)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing workload URI template")
	}
	return &sender{
		log:    logger.New("slack"),
		config: config,

		resourceURITmpl: tmpl,
	}, nil
}

func (s *sender) Compose(event fluxevent.Event) slack.Message {
	blocks := []slack.Block{}
	switch event.Type {
	case fluxevent.EventAutoRelease:
		metadata := event.Metadata.(*fluxevent.AutoReleaseEventMetadata)
		imageIDs := metadata.Result.ChangedImages()
		if len(imageIDs) == 0 {
			imageIDs = []string{"<no image>"}
		}

		// get rid of repository URLs from image ID if we can
		formattedImageIDs := make([]string, len(imageIDs))
		for i, imageID := range imageIDs {
			formattedImageIDs[i] = s.formatImageID(imageID)
		}

		headerTxt := fmt.Sprintf(
			"Automatically released %s",
			strings.Join(formattedImageIDs, ", "),
		)
		blocks = append(blocks, s.formatHeader(headerTxt))

	case fluxevent.EventSync:
		metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
		commitCount := len(metadata.Commits)

		headerTxt := fmt.Sprintf(
			"Synced %d commits to %s",
			len(metadata.Commits),
			s.config.ClusterName,
		)
		blocks = append(blocks, s.formatHeader(headerTxt))

		if commitCount > 0 {
			commits := make([]string, commitCount)
			for i, commit := range metadata.Commits {
				uri := s.getCommitURI(commit.Revision)
				commits[i] = fmt.Sprintf("•  `<%s|%s>` - %s", uri, shortRevision(commit.Revision), commit.Message)
			}
			blocks = append(blocks, headingBlock("Commits")...)
			blocks = append(blocks, textBlock(strings.Join(commits, "\n")))
		}

	default:
		blocks = append(blocks, s.formatHeader(event.String()))
	}

	// common: affected workloads
	affectedWorkloads := make([]string, len(event.ServiceIDs))
	for _, serviceID := range event.ServiceIDs {
		namespace, kind, name := serviceID.Components()
		txt := fmt.Sprintf("•  <%s|%s/%s> in _%s_\n", s.getResourceURI(serviceID), kind, name, namespace)

		affectedWorkloads = append(affectedWorkloads, txt)
	}
	blocks = append(blocks, headingBlock("Affected Workloads")...)
	blocks = append(blocks, textBlock(strings.Join(affectedWorkloads, "\n")))

	msg := slack.NewBlockMessage(blocks...)
	msg.Channel = s.config.SlackChannel
	msg.Username = s.config.SlackUserName
	return msg
}

func (s *sender) formatHeader(text string) slack.Block {
	txt := fmt.Sprintf("*<%s|%s%s>", s.config.VCSRootURL, text, s.config.MessagePostfix)
	return textBlock(txt)
}

func (s *sender) getCommitURI(revision string) string {
	// TODO: customizable commit template
	p := fmt.Sprintf("/commit/%s", revision)
	return path.Join(s.config.VCSRootURL, p)
}

func (s *sender) formatImageID(imageID string) string {
	for _, repoURLToSkip := range s.config.OmittedRepositoryURL {
		if strings.Contains(imageID, repoURLToSkip) {
			return strings.ReplaceAll(imageID, repoURLToSkip, "")
		}
	}
	return imageID
}

func (s *sender) Send(msg slack.Message) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "marshal failed")
	}

	resp, err := http.DefaultClient.Post(
		s.config.SlackWebhookURL,
		"application/json",
		bytes.NewReader(raw),
	)
	if err != nil {
		return errors.Wrap(err, "failed to post webhook")
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return errors.Errorf(
			"API call quota exceed: you need to retry after %v",
			time.Duration(retry)*time.Second,
		)
	}

	// Slack seems to send an HTML body along with 5xx error codes. Don't parse it.
	if resp.StatusCode != http.StatusOK {
		body, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		s.log.Error("Slack returned HTTP {}:\n{}", resp.StatusCode, body)
		return errors.Errorf("slack returned HTTP %d", resp.StatusCode)
	}
	return nil
}

func headingBlock(text string) []slack.Block {
	return []slack.Block{
		slack.NewContextBlock("", slack.NewTextBlockObject(slack.MarkdownType, text, false, false)),
		slack.NewDividerBlock(),
	}
}

func textBlock(text string) slack.Block {
	textObj := slack.NewTextBlockObject(slack.MarkdownType, text, false, false)
	return slack.NewSectionBlock(textObj, nil, nil)
}

func shortRevision(rev string) string {
	if len(rev) <= 7 {
		return rev
	}
	return rev[:7]
}
