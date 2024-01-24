# Remotely Controlled Telemetry

This repo is dedicated to client applications running on customer environments and provide an example implementation of remotely adapting the telemetry data that is gather from the client environment. The problem this repo is trying to tackle is as follows:

The OpenTelemetry SDK is to be instantiated as early as possible within the application where after the instantiation it cannot be restarted with a different configuration. This would mean to restart the entire application.

Instead, an OpenTelemetry collector can be attached to the application which can be started and stopped on demand which allows one to spin it up with different configurations instantly.

The repo shows an example use-case of how this can be achieved.

## Prerequisities

- Golang `1.21`
- New Relic account

## Introduction

Let's say, you are a streaming company. You provide your customers with an application which they run on their own personal computers to connect to a streaming server in your company. From time to time, your customers complain that the streaming is too slow and of course you would like to know why. You implement some debug logs into the client application so that whenever a customer complains, you can troubleshoot the root cause easily.

Everything works fine, BUT... You are sending a ton of debug logs to your observability backend and you overload the network bandwidth. Wouldn't it be great if you could actually turn on and off the sending of these debug logs on demand? Moreover, wouldn't it be great if you could know that there is an issue with a customer application before the customer starts complaining?

## Setup

Let's describe the example environment a bit more detailed.

- You have a customer and he runs your `client` application on his own PC.
- You have a `server` in your company which the `client` application is connecting to.
- You have engineers, SREs etc. who would like to troubleshoot the streaming delay to improve your product.

Your SREs would like to

- know when there is a delay happening in the `client` application.
- turn on the debug logs on the `client` application when that happens.

Here is how you can tackle it:
![setup](/media/setup.png)

### Server

1. Create an HTTP server which can be accessed by your SREs.
2. Create a web socket server which is meant for the `client` application.

### Client

1. Establish a web socket client connection with the `server`.
2. Run the OpenTelemetry collector with the necessary configuration depending on what is told by the `server`.

## Run the environment

### Preparation

Clone the repo:

```shell
git clone git@github.com:utr1903/remotely-controlled-telemetry.git
```

The OpenTelemetry collector binary is not in the repository. Download it:

```shell
mkdir ./apps/client/bin
cd ./apps/client/bin
curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.92.0/otelcol-contrib_0.92.0_darwin_amd64.tar.gz
tar -xvf otelcol-contrib_0.92.0_darwin_amd64.tar.gz
```

### Server

Run the server:

```shell
cd ./apps/server
go run main.go
```

The web socket server will listen to the localhost on the port `8080` and the HTTP server on the port `8081`.

- HTTP server is for the SRE to remotely configure the OpenTelemetry collector of the `client`.
- Web socket server is responsible for delivering the SRE request to the `client`.

### Client

Run the client:

```shell
cd ./apps/client
export OTEL_SERVICE_NAME=client; export NEWRELIC_LICENSE_KEY=<YOUR_LICENSE_KEY>; go run main.go
```

The web socket client will connect to the localhost on the port `8080` and the HTTP server will listen to the localhost on the port `8082`.

- Web socket client is responsible for receiving the SRE request from the `server`.
- HTTP server is for the you to cause a delay in the application for demonstration purposes.

## Monitoring

The client application is already instrumented with OpenTelemetry and it records the latency of a dummy application with the metric `application.latency` which is a histogram. This metric is being sent to the OpenTelemetry collector which then pushes it to the New Relic.

You can query it as follows:

```
FROM Metric SELECT average(application.latency) WHERE service.name = 'client' TIMESERIES SINCE 10 minutes ago
```

Moreover, the `client` application is logging everything to a log file which is tailed by the OpenTelemetry collector that sends the logs (`>INFO`) to New Relic.

You can query the logs as follows:

```
FROM Log SELECT * WHERE service.name = 'client' SINCE 10 minutes ago
```

At first, you might not see any logs or only `INFO` ones.

### Simulating error

To increase the application latency from ~1 second to ~3 seconds, do the following:

```shell
curl -X POST --data '{"duration":"3"}' "http://localhost:8082/latency"
```

This will make an HTTP request to the HTTP server of the `client` and increase the latency eventually. You can see the increase:

```
FROM Metric SELECT average(application.latency) WHERE service.name = 'client' TIMESERIES SINCE 10 minutes ago
```

Now, the SRE knows that there is something bad going on and the needs the debug logs. He does the following:

```shell
curl -X POST "http://localhost:8081/control"
```

This will make an HTTP request to the HTTP server of the `server`. The server passes the request to the web socket and it will eventually trigger the restart of the OpenTelemetry collector of the `client` with debug logs enabled.

Now check the logs again:

```
FROM Log SELECT * WHERE service.name = 'client' SINCE 10 minutes ago
```

You will now see the debug logs.

### Going back to default

The SRE has troubleshooted the issue and everything is now settled. The debug logs are no longer needed. Do the following:

```shell
curl -X POST --data '{"duration":"1"}' "http://localhost:8082/latency"
```

This will reduce the latency back to ~1 second (as if you've solved the issue).

```shell
curl -X DELETE "http://localhost:8081/control"
```

This will make an HTTP request to the HTTP server of the `server`. The server passes the request to the web socket and it will eventually trigger the restart of the OpenTelemetry collector of the `client` with debug logs disabled.
