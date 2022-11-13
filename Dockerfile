FROM golang:1.19.2-alpine3.16 AS base

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go mod download

FROM base as build

RUN go build -o deconz-exporter -tags prod main.go 

FROM scratch as prod

COPY --from=build /build/deconz-exporter /

EXPOSE 8080

ENTRYPOINT ["./deconz-exporter"]