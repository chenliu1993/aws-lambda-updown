
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 go build -o main main.go


.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go


.PHONY: package-linux
package-linux:
	zip function.zip main

.PHONY: package-windows
package-windows:
	go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip && \
		build-lambda-zip -output main.zip main
    
.PHONY: all-linux
all-linux: build-linux package-linux

.PHONY: all-windows
all-windows: build-windows package-windows