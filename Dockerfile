FROM golang:latest
RUN apt-get update

RUN apt-get install -y gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    libgstreamer-plugins-base1.0-dev

WORKDIR $GOPATH/src/github.com/silently/gstreamer-transform

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

CMD ./gstreamer-transform