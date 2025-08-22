FROM golang:1.23.10-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN go build -o /app/fuck-u-code ./cmd/fuck-u-code

FROM alpine:3.20.3
WORKDIR /bin
COPY --from=builder /app/fuck-u-code .

RUN adduser -D nonroot
USER nonroot

ENTRYPOINT ["/bin/fuck-u-code"]
