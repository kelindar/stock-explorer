FROM golang:1.7
MAINTAINER Roman Atachiants "roman@misakai.com"

# copy the local package files to the container's workspace.
ADD . /go/src/github.com/kelindar/stock-explorer
UN go install github.com/kelindar/stock-explorer

# set the entry point
ENTRYPOINT /go/bin/stock-explorer