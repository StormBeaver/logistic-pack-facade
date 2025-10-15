# Builder
FROM golang:1.24.1-alpine AS builder

ARG GITHUB_PATH=github.com/StormBeaver/logistic-pack-facade

WORKDIR /home/${GITHUB_PATH}

RUN apk add --update make git
COPY Makefile Makefile
COPY . .
RUN make go-build

# facade

FROM alpine:latest AS facade

ARG GITHUB_PATH=github.com/StormBeaver/logistic-pack-facade

LABEL org.opencontainers.image.source=https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/facade .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .

RUN chown root:root facade

CMD ["./facade"]
