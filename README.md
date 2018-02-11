# Simple Load Balancer [![CircleCI](https://circleci.com/gh/JackyChiu/slb.svg?style=svg)](https://circleci.com/gh/JackyChiu/slb) [![Go Report Card](https://goreportcard.com/badge/github.com/JackyChiu/slb)](https://goreportcard.com/report/github.com/JackyChiu/slb)

slb is a simple single node HTTP load balancer

Available strategies are: 
- [Round Robin](round_robin.go)
- [Least Busy](least_busy.go)

These strategies for the node pools are also concurrent safe.

## Demo
### Round Robin
```bash
bin/slb -strategy round_robin -config config.json
```
![](.github/round_robin.gif)

### Least Busy
```bash
bin/slb -strategy least_busy -config config.json
```
![](.github/least_busy.gif)

The pool was tested by a requester that produces tasks at a random interval, faster than the time it takes to complete the requests itself. 
The tasks also take an inconsistent amount of time to finish.
In this example the requests were just sleeps.

The standard deviation measures how spread out the pending requests are across workers, giving us an idea of how well the pool is distributing work.
As time goes by the average load goes up the standard deviation remains about the same!

**Basically the smaller the `Std Dev` the better work is being distributed!**

It's worth noting that because the requests aren't equal in the amount of work, the round robin balancer's standard deviation starts to climb.
While the least busy strategy keeps the amount of pending tasks in account and results in a better distribution of requests. 

## Try The Demo
```bash
# clone the project
go get github.com/JackyChiu/slb
cd $GOPATH/src/github.com/JackyChiu/slb

# make project
make all

# start local servers (only for demo)
bin/servers -config config.json

# start balancer server
bin/slb -strategy <round_robin or least_busy> -config config.json

# start the requester (only for demo)
bin/requester
```

A json file is used to specify hosts that are in the node pool, and the port that the load balancer to run on.
```json
{
  "port": 8000,
  "hosts": [
    "cool-app01.mydomain.com",
    "cool-app02.mydomain.com",
    "cool-app03.mydomain.com",
    "cool-app04.mydomain.com"
  ]
}
```
