FROM golang:latest

RUN mkdir -p $GOPATH/src/golang.org/x/ \
    && cd $GOPATH/src/golang.org/x/ \
    && git clone https://github.com/golang/crypto.git \
    && git clone https://github.com/golang/net.git \
    && go get github.com/miekg/dns \
    && go get github.com/op/go-logging \
    && go get github.com/sparrc/go-ping

WORKDIR /dnssecure/src

CMD ["go", "build", "-o", "/dnssecure/bin/dnssecure.b", "./"]
