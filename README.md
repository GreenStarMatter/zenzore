# Zenzore
This project will create mock sensor data and integrate with a GCP Pub/Sub Topic.


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
