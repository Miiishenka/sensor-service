# syntax=docker/dockerfile:1

FROM golang:1.20 as builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o sensor-app homework/cmd/server

FROM alpine:edge

WORKDIR /

COPY --from=builder /sensor-app .

EXPOSE 8080

CMD ["/sensor-app"]