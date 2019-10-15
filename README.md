fluxslack
=====================

Notifies flux events into your own Slack channel using [fluxcloud](https://github.com/justinbarrick/fluxcloud).


## Set Up

You need both [flux](https://github.com/fluxcd/flux) and [fluxcloud](https://github.com/justinbarrick/fluxcloud) to be configured.
For details, please refer their documentation.

### Configuring flux-slack-alert into your cluster

You may need to deploy fluxslack with fluxcloud, with setting up the following configurations.

### Environment Variables

| Name | Description | Required | Default |
|------|-------------|----------|---------|
| **SLACK_WEBHOOK_URL** | Slack Webhook Service URL. | O | |
| **CLUSTER_NAME** | Name of the Kubernetes cluster. | O | |
| **VCS_ROOT_URL** | URL of the GitOps repository. | O | |
| SLACK_CHANNEL | Slack channel ID or name for the alert. | | (Webhook Settings) |
| SLACK_USER_NAME | Displayed user name of the bot. | | (Webhook Settings) |
| MESSAGE_POSTFIX | A message appended to the header message. (e.g. you can tag someone) | | (Webhook Settings) |
| WORKLOAD_URI_TEMPLATE | A URI template for linking workload details | | [kubernetes-dashboard](https://github.com/kubernetes/dashboard) URL, Supposing that it is deployed on `kubernetes-dashboard` namespace with proxy running on `localhost:8001` |
| OMITTED_REPOSITORY_URL | A comma-separated list for known container repository URLs which can be skipped in alert (e.g. Amazon ECR base URL) | |  |

### Configuring fluxcloud

You need to set up following environment variables to fluxcloud:

 - `EXPORTER_TYPE`: `webhook`
 - `EXPORTER_URL`: `http://<ENDPOINT>/v1/webhook`

You may need to delete your Slack integrations from `fluxcloud` if you have configured before.


## License
The fluxslack project is licensed under the MIT License, also included in our repository in the `LICENSE` file.
