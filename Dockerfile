# syntax=docker/dockerfile:1.21
FROM golang:1.26 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o midnight-lyrics-bot

FROM alpine:latest

RUN apk add --no-cache tzdata

COPY --from=build /build/midnight-lyrics-bot ./
COPY --parents albums/ fonts/ images/ ./

ENTRYPOINT [ "./midnight-lyrics-bot" ]
