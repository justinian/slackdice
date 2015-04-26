# slackdice

Slackdice is a golang webservice to implement Slack /roll and /rollprivate commands,
taking any dice argument supported by my [dice](https://github.com/justinian/slackdice)
library. You can set these commands to be anything you'd like when you configure the
integration on the Slack side. Your Slack team will need incoming webhooks configured.
See [Incoming Webhooks](https://my.slack.com/services/new/incoming-webhook).

## Running the service

The service is most easily run as a docker container. The only configuration necessary
is your incoming webhook integration URL.

```bash
docker run -d -e SLACKDICE_SLACK_URL="<your incoming webhook URL>" -p 8000:8000 --name=slackdice justinian/slackdice
```

## Installing into Slack

Just add a slash command integration pointing at your service. Slackdice supports two
urls, `http://your.service:8000/roll` and `http://your.service:8000/roll/private`. In
my slack team, I've set up the `/roll` and `/rollprivate` commands for these, respectively.
