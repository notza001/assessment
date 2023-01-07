# syntax=docker/dockerfile:1

# build
FROM golang:1.19 AS Builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

# deploy
FROM alpine:latest AS Prod
WORKDIR /
COPY --from=Builder /build/app .
CMD ["./app"]