FROM golang:1.7
MAINTAINER Roman Atachiants "roman@misakai.com"

ADD . /go/src/github.com/kelindar/stock-explorer
RUN go install github.com/kelindar/stock-explorer

ENTRYPOINT /go/bin/stock-explorer