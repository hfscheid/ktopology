FROM golang:1.23

WORKDIR /usr/src/app
RUN mkdir ktprom

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY simulator/svcsim/go.mod simulator/svcsim/go.sum ./
COPY ktprom/ktprom.go ktprom/
COPY simulator/svcsim/ .
RUN go mod tidy 
# RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/app ./.

EXPOSE 80

CMD ["app"]

