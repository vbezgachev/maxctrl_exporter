FROM golang:1.12

RUN go get -u github.com/golang/dep/cmd/dep
COPY . /go/src/github.com/maxctrl_exporter

RUN cd /go/src/github.com/maxctrl_exporter && \
    make build

RUN mv /go/src/github.com/maxctrl_exporter/maxctrl_exporter /bin/maxctrl_exporter

ENTRYPOINT  ["/bin/maxctrl_exporter"]
