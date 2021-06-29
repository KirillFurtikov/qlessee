tests:
	ginkgo pkg/qless/tests

build:
	go build -ldflags '-s -w' .

run:
	QLESS_GO_REDIS_LIST="foobar=redis://127.0.0.1:6379" go run .

compile:
    GOOS=freebsd GOARCH=386 go build -ldflags '-s -w' -o bin/qlessee-freebsd-386 .
    # Linux
    GOOS=linux GOARCH=386 go build -ldflags '-s -w' -o bin/qlessee-linux-386 .
    # Windows
    GOOS=windows GOARCH=386 go build -ldflags '-s -w' -o bin/qlessee-windows-386 .
        # 64-Bit
    # FreeBDS
    GOOS=freebsd GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-freebsd-amd64 .
    # MacOS
    GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-darwin-amd64 .
    # Linux
    GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-linux-amd64 .
    # Windows
    GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o bin/qlessee-windows-amd64 .
