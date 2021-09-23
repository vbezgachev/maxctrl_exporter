FROM golang:1.17

COPY . /go/src/github.com/maxctrl_exporter

RUN cd /go/src/github.com/maxctrl_exporter \
    && go build ./... \
    && mv /go/src/github.com/maxctrl_exporter/maxctrl_exporter /bin/maxctrl_exporter

ENTRYPOINT  ["/bin/maxctrl_exporter"]
