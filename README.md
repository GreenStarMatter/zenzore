# Zenzore
> Version 0.1.0
This project will create mock sensor data and integrate with a GCP Pub/Sub Topic.


Zenzore is a CLI which mocks devices containing sensors and passes sampleddata out in a bulk message.  The app will be pseudo-simulated meaning that there will be some mathematical constructs which bind devices, sensors, and samples; however, the app will not be a strict physical simulation.  The focus is on the transfer of information to create the basis of a data pipeline.

## App Command Structure
- zenzore
    - run (starts all zenzore zyztems)
    - update (kick into a navigator which allows you to traverse data structure and update)
    - diagnostics (gets running state of zyztems with basic stats)

## Top Level Folder Structure

- root
    - main.go
    - cmd
        - update
        - diagnostics
        - run
    - internal
        - zyztem
        - navigator
        - message

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

**GitHub Pipelines**

`Build App`

uses GitHub actions to get appropriate Go version and then builds zenzore

`Test Application`

runs test, coverage, and report on built zenzore program then archives report


The main architecture will be:
1. A signal generator which creates a signal with some noise
1. A user input layer which will take CLI arguments to make a signal with metadata
1. An authentication and transport layer which will output messages to the Pub/Sub


### Signal Generator
This should create a stastical model for a process which is then sampled from:
1. Create model or waveform for part running
1. Allow user to input parameters for model and noise selection
1. Allow user to sample from different points on the model

### User Input
This should allow the user to feed different information to the program for ease of use:
1. Allow signal attributes (potentially seed) to be fed to signal generator
1. Allow user to indicate amount of sensors being setup, duration, and mode (continuous or discrete)
1. Allow user to assign additional metadata to sensors (common parent devices, delivery paths, etc.)

### Authentication and Messaging
This layer should authenticate with GCP and transfer the messages from the sensors to the topic:
1. Should make sure that system is signed into GCP
1. Should send message using GCP Go libraries
1. Should make log failed message sends (not a failure)
