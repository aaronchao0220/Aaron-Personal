
# usage-telemetry-publisher

[![Maintainability](https://qlty.sh/badges/0ce04001-48a7-4916-b0cf-bf84263c3af9/maintainability.svg)](https://qlty.sh/gh/qlik-trial/projects/usage-telemetry-publisher)
[![Code Coverage](https://qlty.sh/badges/0ce04001-48a7-4916-b0cf-bf84263c3af9/coverage.svg)](https://qlty.sh/gh/qlik-trial/projects/usage-telemetry-publisher)
## Overview

The `usage-telemetry-publisher` This service is responsible for providing access for 3rd party tools to have access to usage telemetry data. It will mask PII and is the mean to publish regional data to a central region use case.

## Features


## Configuration

The service can be configured using a Helm chart. Below are some key configuration options available in the [`values.yaml`](./manifests/chart/usage-telemetry-publisher/values.yaml) file.


## Development

### Building the Project

To build the project, run the following command:

```sh
make build
```

### Running Tests

To run the tests, use the following command:

```sh
make test
```

### Linting the Code

To lint the code, use the following command:

```sh
make lint
```

### Debugging in .vscode

create a launch.json:
```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
          "name": "Launch usage-telemetry-publisher",
          "type": "go",
          "request": "launch",
          "mode": "auto",
          "program": "${workspaceFolder}/cmd/main",
            "env": {
              "LOG_LEVEL": "debug",
              "OTLP_AGENT_HOST": "localhost",
              "OTLP_AGENT_PORT": "4317",
              "AUTH_ENABLED": "false",
              "SOLACE_MESSAGE_VPN": "default",
              "SOLACE_URI": "tcp://localhost:55554",
              "LAUNCHDARKLY_STREAM_URI": "http://localhost:8082/relay",
              "LAUNCHDARKLY_ENABLED": "true",
              "SOLACE_CHANNELS": "ui-events.analytics",
              "INTERMEDIATE_STORAGE_ENABLED": "true",
              "MESSAGING_ENABLED": "true",
              "EVENTS_FILE_PATH": "/etc/config/test/events.yaml",
              "LAUNCHDARKLY_SDK_KEY": "key"
            },
          "args": [],
          "showLog": true
        }

    ]
}
```
And run 
```
make start-dependencies
```
before running debug mode.

Or for running all in docker containers: (includes a reset of all docker containers and volumes)
```
make start-nobuild
```

### Feature Flags



## Contact

For any questions or support, please reach out to the owning team `telemetry-data-services` on Slack channels:

- Discussion: `#tds-public`
- Bot: `#tds-bot`

