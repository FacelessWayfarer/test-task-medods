FROM golang:1.24-alpine AS builder

WORKDIR /src/app

COPY ["go.mod","go.sum","./"]

RUN go mod download

COPY . .

RUN go test -c -o test API_test.go

CMD ["./test"]