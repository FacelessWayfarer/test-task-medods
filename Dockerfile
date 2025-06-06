FROM golang:1.24-alpine AS builder

WORKDIR /src/app

COPY ["go.mod","go.sum","./"]

RUN go mod download

COPY . .

# Download install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Build swagger using swag
RUN swag init -d cmd --pdl 3

RUN go build -o ./bin/app cmd/main.go

FROM alpine:latest AS runner

#Curl is used for testing
RUN apk --no-cache add curl

#Copy bin
COPY --from=builder /src/app/bin/app /

#Copy config
COPY config/local.yaml config/local.yaml

#Copy swagger
COPY docs docs

#Copy migrations
COPY internal/migrations internal/migrations

CMD [ "/app" ]

