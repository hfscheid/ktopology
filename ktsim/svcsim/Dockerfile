FROM golang:1.23

WORKDIR /usr/src/app
# RUN mkdir ktprom

# COPY simulator/svcsim/go.mod simulator/svcsim/go.sum ./
# COPY ktprom/ktprom.go ktprom/
COPY . .
RUN go mod tidy 

RUN go build -v -o /usr/local/bin/app ./.

EXPOSE 80

CMD ["app"]
