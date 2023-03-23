FROM golang:1.20-alpine as builder

WORKDIR /test-server

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags -a -o bin/test-server main.go

FROM ubuntu

RUN adduser --system --group app \
    && apt-get -y update && apt-get install -y curl

WORKDIR /app

ARG TEMPORAL_TEST_SERVER_VERSION=1.18.2

RUN curl -LO https://github.com/temporalio/sdk-java/releases/download/v${TEMPORAL_TEST_SERVER_VERSION}/temporal-test-server_${TEMPORAL_TEST_SERVER_VERSION}_linux_amd64.tar.gz \
    && tar -xzf temporal-test-server_${TEMPORAL_TEST_SERVER_VERSION}_linux_amd64.tar.gz \
    && mv temporal-test-server_${TEMPORAL_TEST_SERVER_VERSION}_linux_amd64/temporal-test-server /app/temporal-test-server \
    && rm -rf temporal-test-server_${TEMPORAL_TEST_SERVER_VERSION}_linux_amd64.tar.gz \
    && chmod +x /app/temporal-test-server

COPY --from=builder /test-server/bin/test-server .

RUN chown -R app:app ./

USER app

EXPOSE 7233 1323

CMD ["./test-server"]
