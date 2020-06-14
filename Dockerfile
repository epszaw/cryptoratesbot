FROM golang:1.13

WORKDIR .
COPY . /go/src/lamartire/cryptoratesbot
WORKDIR /go/src/lamartire/cryptoratesbot
RUN go mod download
RUN go mod vendor
RUN go build main.go
CMD ["./main"]
