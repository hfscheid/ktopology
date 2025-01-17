FROM golang:1.23

WORKDIR /usr/src/app
COPY test/main.go .
RUN mkdir poddiscovery
COPY poddiscovery.go ./poddiscovery
RUN go mod init testpoddiscovery && go mod tidy

RUN go build -v -o /usr/local/bin/app ./.

CMD ["app"]
