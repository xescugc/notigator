# build stage
FROM golang:1.12 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# final stage
FROM alpine

RUN apk --update add ca-certificates

COPY --from=builder /app/notigator /app/

ENTRYPOINT ["/app/notigator"]
