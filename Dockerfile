FROM scratch

ADD ["bin/dnssecure-bin", "/"]

CMD ["/dnssecure-bin"]