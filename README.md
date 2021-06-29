# QLessee

## _Terminal user interface for Qless_

It was inspired by [k9s](https://github.com/derailed/k9s) for Kubernetes and [docui](https://github.com/skanehira/docui) for Docker. Based on [tview](https://github.com/rivo/tview)

## Features

- Show queues statuses
- Allow to pause and continue queues
- Failed jobs navigation and inspection
- Inspect job details

## Requirements

Environment variable like:

```sh
QLESS_GO_REDIS_LIST="main=redis://217.0.0.1:6379/0,external=redis://217.0.0.1:6380/0"
```

Do not use it for production environments cuz alpha, ok?

For running tests, build, run - look at Makefile

```makefile
tests:
    ginkgo pkg/qless/tests

build:
    go build .

run:
    QLESS_GO_REDIS_LIST="foobar=redis://127.0.0.1:6379" go run .

compile:
    GOOS=freebsd GOARCH=386 go build -ldflags '-s -w' -o bin/qlessee-freebsd-386 .
    GOOS=linux GOARCH=386 go build -ldflags '-s -w' -o bin/qlessee-linux-386 .
    GOOS=windows GOARCH=386 go build -ldflags '-s -w' -o bin/qlessee-windows-386 .
    GOOS=freebsd GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-freebsd-amd64 .
    GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-darwin-amd64 .
    GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-linux-amd64 .
    GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-windows-amd64 .
```
