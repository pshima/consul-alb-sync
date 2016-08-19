# consul-alb-sync

A proof of concept to sync AWS Application Load Balancer (ALB) target groups with consul service discovery.

## Usage

Create a consul k/v structure containing 2 keys
* Enabled (set to true)
* ServiceName (set to the name of the consul service)

Right now the Target group must be created but would be easy to add in functionality to auto create target groups.

This will sync the consul service with a target group of the same name, matching host and port to instance ID and registering or deregistering hosts behind the target group.

It only runs once, but in theory would just run in a loop.

