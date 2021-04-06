APPNAME = cardbill
BIN = $(GOPATH)/bin
GOCMD = /usr/local/go/bin/go
GOBUILD = $(GOCMD) build
GOINSTALL = $(GOCMD) install
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GORUN = $(GOCMD) run
BINARY_UNIX = $(BIN)_unix
PID = .pid
HUB_ADDR = hub.kpaas.nsini.com
TAG = v0.0.01-test
NAMESPACE = app
PWD = $(shell pwd)
GOPROXY = https://goproxy.cn

start:
	$(GOBUILD) -v
	$(BIN)/$(APPNAME) & echo $$! > $(PID)

restart:
	@echo restart the app...
	@kill `cat $(PID)` || true
	$(BIN)/$(APPNAME) & echo $$! > $(PID)

stop:
	@kill `cat $(PID)` || true

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PWD) $(HUB_ADDR)

build:
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o cardbill -v ./cmd/server/service.go
	docker build --rm -t $(APPNAME):$(TAG) .

docker-run:
	docker run -it --rm -p 8080:8080 -v $(PWD)/app.cfg:/etc/cardbill/app.cfg $(APPNAME):$(TAG)

push:
	docker image tag $(APPNAME):$(TAG) $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(TAG)
	docker push $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(TAG)

run:
	GO111MODULE=on GOPROXY=$(GOPROXY) $(GORUN) ./main.go -http-addr :8080 -config-file ./app.cfg

client-init:
	GO111MODULE=on $(GORUN) ./cmd/client/client.go -config-file ./app.cfg

get:
	GO111MODULE=on $(GOGET) $(shell v='$(FULL_VERSION)'; echo "$${v%.*}")
