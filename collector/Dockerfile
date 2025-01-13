FROM golang:1.23

WORKDIR /usr/src/app
RUN mkdir -p pkg/ktprom

# copy collector/go.mod collector/go.sum ./
COPY ktprom/ktprom.go ktprom/
COPY collector/ .
RUN go mod tidy

RUN go build -v -o /usr/local/bin/app ./.

EXPOSE 80

CMD ["app"]
