Create configs from templates and environment variables then execute command.

Useful for configuring applications in 12 factor environments (eg. Kubernetes) where the applications don't use environment variables for config.

Example for logstash provided

```
. test.source
go build templator.go
./templator cat logstash.conf
```
