FROM golang:1.23.2

WORKDIR /usr/app

COPY go.mod go.sum ./

RUN go mod download

COPY ./btc.go ./counter.go ./main.go ./telegram.go ./

RUN go build -o bruter .

CMD ["./bruter"]
