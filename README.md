# Zenzore
> Version 1.1.0

This project will create mock sensor data and integrate with a GCP Pub/Sub Topic.


Zenzore is a CLI which mocks devices containing sensors and passes sampled data out in a bulk message.  The app will be pseudo-simulated meaning that there will be some mathematical constructs which bind devices, sensors, and samples; however, the app will not be a strict physical simulation.  The focus is on the transfer of information to create the basis of a data pipeline.

## Generic Architecture

**Backend**
1. A CLI front-end for creating zyztems, managing their state, and persisting
1. A server back-end allowing for in memory persistence of state and organizing user interfaces with the app logic
1. An authentication and transport layer which will output messages to Pub/Sub using ADC

**GCP Data Pipeline**
1. A Pub/Sub ingest topic receiving messages from the Zenzore backend
1. A BigQuery subscription writing messages directly to a raw events table using the table schema
1. A dead-letter topic and Cloud Storage sink for capturing failed messages
1. All infrastructure managed via Terraform in the `infra/` directory

## Setup Data Pipeline

### 1. Configure environment variables
Copy the example env file and fill in your own values:
```bash
cp infra/EXAMPLE infra/.env.tf
```
Edit `infra/.env.tf` with your GCP project ID and resource names, then source it:
```bash
source infra/.env.tf
```
This must be repeated in each new terminal session before running Terraform commands.

### 2. Run Terraform
Initialize, plan, and apply the infrastructure:
```bash
make tf-init
make tf-plan
make tf-apply
```
To tear down all managed infrastructure:
```bash
make tf-destroy
```

## App Command Structure

- zenzore
    - run (starts a zenzore server that allows for updates)
    - nav (allows for easy menu-like navigation of an existing server, purely for exploring)
    - diagnostics (gets running state of zyztems with basic stats)

## Top Level Folder Structure

- root
    - main.go
    - cmd
        - run
        - nav
        - diagnostics (NOT IMPLEMENTED YET)
    - internal
        - zyztem
        - navigator
        - message
        - server
        - appdata (deprecated, but likely to be replaced by db persistence)
    - infra


## Environment Variables

Zenzore reads its configuration from environment variables rather than a config file. The server and the GCP message pipeline each require their own set, listed below.

### Server (`zenzore run`)

| Variable | Required | Description |
|---|---|---|
| `ZENZORE_PORT` | Yes | Port the root server listens on. No default is provided — `run` will fail with an explicit error if this is unset, so a port is never silently assumed. |

### GCP Pub/Sub (`message` package, used by `/zyztems/send`)

| Variable | Required | Description |
|---|---|---|
| `ZENZOREPROJECTID` | Yes | The GCP project ID the Pub/Sub topic belongs to. |
| `ZENZORETOPICID` | Yes | The Pub/Sub topic name messages are published to. |

Both are required only when triggering a send (e.g. hitting `/zyztems/send`); they are not needed to run `nav` or `diagnostics` against an already-running server.

### Connecting to a running server (`zenzore nav`)

`nav` does not read any server-related environment variable automatically — it has no way to know which server you mean unless you tell it. Pass the server's address explicitly:

```bash
zenzore nav --server http://localhost:8080
```

## Known Limitations

**Concurrency is only partially handled.** The server's in-memory registry (the map tracking all zyztems) is protected by a mutex, so concurrent requests against *different* zyztems, or concurrent create/list/remove calls, are safe. However, once a specific `Zyztem` is retrieved from the registry, mutations to its own state (adding a device, adding a sensor, sampling) are **not yet locked at that level**. Two concurrent requests targeting the *same* zyztem (for example, two simultaneous calls to add a device to the same zyztem) can race.

In practice, this means Zenzore is currently safe for:
- Multiple zyztems being created, listed, or removed concurrently
- A single client interacting with the server at a time

It is **not yet safe** for:
- Multiple clients concurrently mutating the *same* zyztem (adding devices/sensors, sampling) at the same time

Per-zyztem locking is planned but not yet implemented. Until then, avoid issuing concurrent mutating requests against the same zyztem ID.

## Development
This project uses a `Makefile` to simplify common tasks.  These are used within a GitHub Pipeline controlled by a yaml file for CI.

**Command**

`make init-go`

initializes the the Go environment by adding its path to the bash environment

`make build`

builds program binary

`make report`

creates a coverage report from output of an existing machine readable coverage profile

`make test`

runs tests over entire project and creates machine readable coverage profile

`make coverage`

runs a pass fail gate by measure coverage reater than 80 percent.

`make tf-init`
initializes Terraform and downloads the GCP provider in the `infra/` directory

`make tf-plan`
shows a dry-run of infrastructure changes without applying them

`make tf-apply`
creates or updates all GCP infrastructure defined in `infra/`

`make tf-destroy`
tears down all Terraform-managed GCP infrastructure

**GitHub Pipelines**

`Build App`

uses GitHub actions to get appropriate Go version and then builds zenzore

`Test Application`

runs test, coverage, and report on built zenzore program then archives report
