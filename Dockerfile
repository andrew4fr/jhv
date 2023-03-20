FROM golang:alpine as builder

WORKDIR /builds/

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .

RUN go build -a -o /tmp/treasure ./cmd/treasure-api-server/main.go


FROM alpine

ARG APP_NAME=app

WORKDIR /app/

COPY --from=builder /tmp/treasure /app/treasure

ENTRYPOINT ["./treasure"]

