FROM golang:1.18-alpine
WORKDIR /go/src/app
COPY app/. .

ENV GODEBUG=http2client=0
ENV CGO_ENABLED=0
ENV CC=gcc
RUN go get -d -v ./...
RUN go install -v ./...
RUN which github-pull

from ubuntu:focal

WORKDIR /root/
COPY --from=0 /go/bin/github-pull ./
CMD ["./github-pull"]