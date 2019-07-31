# Build stage
FROM golang:1.12.7-alpine3.10 AS build

ENV GOPATH=/go \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY main.go main.go

RUN go build -a -installsuffix cgo -ldflags="-s -w" -tags="jsoniter" -o /bin/server main.go

# Production stage
FROM scratch

ENV LANG=C.UTF-8

WORKDIR /app

COPY --from=build /bin/server .

COPY direct.txt .

ENTRYPOINT ["./server"]

CMD ["--direct", "direct.txt"]
