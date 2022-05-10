FROM golang:1.11

VOLUME ["/go"]

WORKDIR /go/src/github.com/dnerochain/dnero-eth-rpc/

ENV GOPATH=/go

ENV CGO_ENABLED=1 

ENV GO111MODULE=on

CMD ["/go/src/github.com/dnerochain/dnero-eth-rpc/integration/docker/build/start.sh"]



