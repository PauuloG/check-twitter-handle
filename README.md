# check-twitter-handle
Check availability of a list of twitter handles

Starting December 11th 2019, Twitter is going to deactivate inactive accounts.
See this article for info : [Twitter prepares for huge cull of inactive users](https://www.bbc.co.uk/news/technology-50567751)


This simple command checks for availability of a list of handles and notifies you.
For now the app is only able to notify users on Slack. (Feel free to PR for other notifier gateways)

## Installation

Easiest way to use this is to compile it in a docker container :
```
docker build -t check-twitter-handle:latest -f ./Dockerfile ./
```

## Running paris-pound-check

Then, running the script itself is a matter of populating environment variables in `.env` and starting docker container with environment file mounted

```
docker run --rm -v $(pwd)/.env:/go/bin/.env check-twitter-handle:latest --notifier=slack check
```

The official image built from this repository is also available on Docker hub image registry [https://cloud.docker.com/u/pauulog/repository/docker/pauulog/check-twitter-handle](https://cloud.docker.com/u/pauulog/repository/docker/pauulog/check-twitter-handle)

### Cron usage

In the following example, the program will check every minute if the vehicle has been impounded.
```
* * * * docker run --rm -v $(pwd)/.env:/go/bin/.env docker.io/pauulog/check-twitter-handle:latest --notifier=slack check
```
