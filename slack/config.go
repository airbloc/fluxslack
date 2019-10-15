package slack

type Config struct {
	SlackWebhookURL string `envconfig:"slack_webhook_url" required:"true"`
	ClusterName     string `envconfig:"cluster_name" required:"true"`
	VCSRootURL      string `envconfig:"vcs_root_url" required:"true"`

	SlackChannel   string `envconfig:"slack_channel"`
	SlackUserName  string `envconfig:"slack_user_name"`
	MessagePostfix string `envconfig:"message_postfix"`

	// Base URI for showing kubernetes workloads
	WorkloadURITemplate string `envconfig:"workload_uri_template" default:"http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/"`

	OmittedRepositoryURL []string `envconfig:"omitted_repository_url"`
}
