FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache bash

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o main .

FROM alpine:latest AS runner

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /app/main .
COPY . .


EXPOSE 8080

CMD ["./main"]