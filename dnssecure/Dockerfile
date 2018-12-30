FROM golang:latest

RUN mkdir -p $GOPATH/src/golang.org/x/ \
    && cd $GOPATH/src/golang.org/x/ \
    && git clone https://github.com/golang/net.git \
    && go get github.com/miekg/dns \
    && go get github.com/go-sql-driver/mysql \
    && go get github.com/sparrc/go-ping

ADD ["go-build.sh", "/"]

WORKDIR /
