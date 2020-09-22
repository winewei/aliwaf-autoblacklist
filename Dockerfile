FROM golang:1.14-alpine3.12 AS builder

WORKDIR /srv

COPY . .

RUN go build -o aliwaf-autoblacklist main.go

FROM alpine

WORKDIR /srv

COPY --from=builder /srv/aliwaf-autoblacklist aliwaf-autoblacklist

CMD ["/srv/aliwaf-autoblacklist"]
