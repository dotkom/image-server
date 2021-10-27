FROM golang:1.17 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .
COPY cmd cmd
COPY internal internal

RUN go build -o ./image-server

FROM gcr.io/distroless/base

WORKDIR /usr/image-server

COPY --from=build /build/image-server /image-server

EXPOSE 8080

ENTRYPOINT ["/image-server", "server"]