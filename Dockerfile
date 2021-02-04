FROM golang:stretch AS build

COPY .  /go/src/grafana-wechat
WORKDIR /go/src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/grafana-wechat-alert *.go

FROM alpine

COPY --from=build /go/bin/grafana-wechat-alert /grafana-wechat-alert
ENV DOCKER=chimojiacai
CMD ["/grafana-wechat-alert"]