FROM golang:latest as build-env

ENV GO111MODULE=on
ENV BUILDPATH=gitlab.yixinonline.org/kplcloud/kpaas
ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/go
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM alpine:latest

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        curl \
        && update-ca-certificates 2>/dev/null || true

COPY --from=build-env /go/bin/cardbill /go/bin/cardbill
COPY ./dist /go/bin/dist

WORKDIR /go/bin/
CMD ["/go/bin/cardbill", "-http-addr", ":8080", "-config-file", "/etc/cardbill/app.cfg"]