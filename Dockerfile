FROM debian:11-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /cic-dw

COPY cic-dw .
COPY config.toml .
COPY queries .

CMD ["./cic-dw"]
