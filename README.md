flux-slack-alert
=====================

Notifies flux events into your own Slack channel using [fluxcloud](https://github.com/justinbarrick/fluxcloud).


## Set Up

You need both [flux](https://github.com/fluxcd/flux) and [fluxcloud](https://github.com/justinbarrick/fluxcloud) to be configured.
For details, please refer their documentation.

### Configuring flux-slack-alert into your cluster

### Environment Variables


### Configuring fluxcloud

You need to set up following environment variables to fluxcloud:

 - `EXPORTER_TYPE`: `webhook`
 - `EXPORTER_URL`: `http://flux-slack-alert.<YOUR_NAMESPACE>/v1/webhook`

You may need to delete your Slack integrations from `fluxcloud` if you have configured before.


## License
The airbloc-go project is licensed under the MIT License, also included in our repository in the `LICENSE` file.
