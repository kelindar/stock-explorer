FROM golang:1.7
MAINTAINER Roman Atachiants "roman@misakai.com"

RUN go get -d github.com/kelindar/stock-explorer
RUN go install github.com/kelindar/stock-explorer

# set the entry point
ENTRYPOINT /go/bin/stock-explorer