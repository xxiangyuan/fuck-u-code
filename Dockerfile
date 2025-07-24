FROM golang:1.23.10-alpine

WORKDIR /build

RUN mkdir -p /fuck-u-code

COPY . /fuck-u-code

RUN cd /fuck-u-code && go build -o /bin/fuck-u-code ./cmd/fuck-u-code

ENTRYPOINT ["/bin/fuck-u-code"]