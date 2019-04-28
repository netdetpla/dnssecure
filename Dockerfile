FROM scratch

ADD ["bin/dnssecure.b", "/"]

CMD ["/dnssecure.b"]