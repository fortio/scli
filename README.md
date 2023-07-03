[![Go Reference](https://pkg.go.dev/badge/fortio.org/scli.svg)](https://pkg.go.dev/fortio.org/scli)
[![Go Report Card](https://goreportcard.com/badge/fortio.org/scli)](https://goreportcard.com/report/fortio.org/scli)
[![GitHub Release](https://img.shields.io/github/release/fortio/scli.svg?style=flat)](https://github.com/fortio/scli/releases/)
# scli

Extends [cli](https://github.com/fortio/cli#cli) to server `main()`s .

In addition to flags, usage and help output, arguments validation, `scli` abstracts the repetitive parts of a `main()` to setup a config directory watch for [dynamic flags](https://github.com/fortio/dflag) (configmap in kubernetes cases) and configuration endpoint/UI/api.

You can see real use example in a server like [proxy](https://github.com/fortio/proxy).

## Server Example

Server example [sampleServer](sampleServer/main.go)

Previous style (non structured json log format):
```bash
% go run . -config-dir ./config -config-port 8888 -json-log=false a b
14:50:54 I updater.go:47> Configmap flag value watching on ./config
14:50:54 I updater.go:156> updating loglevel to "verbose\n"
14:50:54 I logger.go:183> Log level is now 1 Verbose (was 2 Info)
14:50:54 I updater.go:97> Now watching . and config
14:50:54 I updater.go:162> Background thread watching config now running
14:50:54 I scli.go:81> Fortio scli dev dflag config server listening on [::]:8888
14:50:54 I scli.go:90> Starting sampleServer dev  go1.20.3 arm64 darwin
14:50:55 I main.go:16> FD count 1s after start : 16
# When visiting the UI
14:51:06 ListFlags: GET / HTTP/1.1 [::1]:64731 () ...
# ...
14:51:15 I main.go:18> FD count 20s later      : 16
14:51:15 I main.go:21> FD count stability check: 16
14:51:15 I main.go:21> FD count stability check: 16
14:51:15 I main.go:21> FD count stability check: 16
14:51:15 I main.go:21> FD count stability check: 16
14:51:15 I main.go:21> FD count stability check: 16
14:51:15 I main.go:27> Running until interrupted (ctrl-c)...
# pkill -int sampleServer
14:51:20 W scli.go:101> Interrupt received.
14:51:20 I main.go:29> Normal exit
% echo $?
0
```

With the flags ui on http://localhost:8888

<img width="716" alt="flags UI" src="https://user-images.githubusercontent.com/3664595/219904547-368a024e-1d6a-4301-a7a9-8882e37f5a90.png">

New default style of logging since 1.5 (JSON for servers):
```bash
$ go run . -config-dir ./config -config-port 8888 a b
```
```json
{"ts":1686609103.447926,"level":"info","file":"updater.go","line":47,"msg":"Configmap flag value watching on ./config"}
{"ts":1686609103.449381,"level":"info","file":"updater.go","line":156,"msg":"updating loglevel to \"verbose\\n\""}
{"ts":1686609103.449406,"level":"info","file":"logger.go","line":224,"msg":"Log level is now 1 Verbose (was 2 Info)"}
{"ts":1686609103.450125,"level":"info","file":"updater.go","line":97,"msg":"Now watching . and config"}
{"ts":1686609103.450240,"level":"info","file":"updater.go","line":162,"msg":"Background thread watching config now running"}
{"ts":1686609103.450523,"level":"info","file":"scli.go","line":87,"msg":"Fortio scli dev dflag config server listening on [::]:8888"}
{"ts":1686609103.450534,"level":"info","file":"scli.go","line":96,"msg":"Starting sampleServer dev  go1.20.5 arm64 darwin"}
{"ts":1686609104.452193,"level":"info","file":"main.go","line":16,"msg":"FD count 1s after start : 14"}
# list flag (curl localhost:8888)
{"ts":1686609330.309960,"level":"info","file":"http_logging.go","line":73,"msg":"ListFlags","method":"GET","url":"/","proto":"HTTP/1.1","remote_addr":"127.0.0.1:60554","header.x-forwarded-proto":"","header.x-forwarded-for":"","user-agent":"curl/8.0.1","header.host":"localhost:8888","header.User-Agent":"curl/8.0.1","header.Accept":"*/*"}
{"ts":1686609124.453697,"level":"info","file":"main.go","line":18,"msg":"FD count 20s later      : 14"}
{"ts":1686609124.454075,"level":"info","file":"main.go","line":21,"msg":"FD count stability check: 14"}
{"ts":1686609124.454411,"level":"info","file":"main.go","line":21,"msg":"FD count stability check: 14"}
{"ts":1686609124.454745,"level":"info","file":"main.go","line":21,"msg":"FD count stability check: 14"}
{"ts":1686609124.455071,"level":"info","file":"main.go","line":21,"msg":"FD count stability check: 14"}
{"ts":1686609124.455462,"level":"info","file":"main.go","line":21,"msg":"FD count stability check: 14"}
{"ts":1686609124.455482,"level":"info","file":"main.go","line":27,"msg":"Running until interrupted (ctrl-c)..."}
# After ^C
{"ts":1686609129.019649,"level":"warn","file":"scli.go","line":107,"msg":"Interrupt received."}
{"ts":1686609129.019703,"level":"info","file":"main.go","line":29,"msg":"Normal exit"}
```

## Additional builtins
(coming from `cli`'s base module)

### buildinfo

e.g

```bash
$ go install fortio.org/cli/sampleServer@latest
go: downloading fortio.org/cli v0.1.0
$ sampleServer buildinfo
0.1.0 h1:M+coC8So/41xSyGiiJ/6RS+XhnshNUslZuaB6H8z9GI= go1.19.6 arm64 darwin
go	go1.19.6
path	fortio.org/cli/sampleServer
mod	fortio.org/cli	v0.1.0	h1:M+coC8So/41xSyGiiJ/6RS+XhnshNUslZuaB6H8z9GI=
dep	fortio.org/dflag	v1.4.1	h1:WDhlHMh3yrQFrvspyN5YEyr8WATdKM2dUJlTxsjCDtI=
dep	fortio.org/fortio	v1.50.1	h1:5FSttAHQsyAsi3dzxDmSByfzDYByrWY/yw53bqOg+Kc=
dep	fortio.org/log	v1.2.2	h1:vs42JjNwiqbMbacittZjJE9+oi72Za6aekML9gKmILg=
dep	fortio.org/version	v1.0.2	h1:8NwxdX58aoeKx7T5xAPO0xlUu1Hpk42nRz5s6e6eKZ0=
dep	github.com/fsnotify/fsnotify	v1.6.0	h1:n+5WquG0fcWoWp6xPWfHdbskMCQaFnG6PfBrh1Ky4HY=
dep	github.com/google/uuid	v1.3.0	h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=
dep	golang.org/x/exp	v0.0.0-20230213192124-5e25df0256eb	h1:PaBZQdo+iSDyHT053FjUCgZQ/9uqVwPOcl7KSWhKn6w=
dep	golang.org/x/net	v0.7.0	h1:rJrUqqhjsgNp7KqAIc25s9pZnjU7TUcSY7HcVZjdn1g=
dep	golang.org/x/sys	v0.5.0	h1:MUK/U/4lj1t1oPg0HfuXDN/Z1wv31ZJ/YcPiGccS4DU=
dep	golang.org/x/text	v0.7.0	h1:4BRB4x83lYWy72KwLD/qYDuTu7q9PjSagHvijDw7cLo=
build	-compiler=gc
build	CGO_ENABLED=1
build	CGO_CFLAGS=
build	CGO_CPPFLAGS=
build	CGO_CXXFLAGS=
build	CGO_LDFLAGS=
build	GOARCH=arm64
build	GOOS=darwin
```

### help
```bash
sampleServer 1.6.0 usage:
	sampleServer [flags] arg1 arg2 [arg3...arg4]
or 1 of the special arguments
	sampleServer {help|version|buildinfo}
flags:
  -config-dir directory
    	Config directory to watch for dynamic flag changes
  -config-port port
    	Config port to open for dynamic flag UI/api
  -logger-file-line
    	Filename and line numbers emitted in JSON logs, use -logger-file-line=false to disable (default true)
  -logger-json
    	Log in JSON format, use -logger-json=false to disable (default true)
  -logger-timestamp
    	Timestamps emitted in JSON logs, use -logger-timestamp=false to disable (default true)
  -loglevel level
    	log level, one of [Debug Verbose Info Warning Error Critical Fatal] (default Info)
  -quiet
    	Quiet mode, sets loglevel to Error (quietly) to reduces the output
```

### Server log diff'ing

When debugging in dev mode the differences between 2 log output, it's convenient to use the following flags

- `-logger-timestamp=false` so the timestamp is removed from the output as that would be different always
- `-logger-file-line=false` so code line numbers don't show as diffs either (if comparing different versions/releases)

### version
Short 'numeric' version (v skipped, useful for docker image tags etc)
```bash
$ sampleServer version
1.6.0
```
