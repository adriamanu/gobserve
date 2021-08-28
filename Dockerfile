FROM golang:1.16 AS builder

WORKDIR /go/goverwatch

COPY go.* ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o goverwatch

ENTRYPOINT [ "/go/goverwatch/goverwatch" ]