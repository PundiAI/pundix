# compile pundix
FROM golang:1.18.2-alpine3.16 as builder

RUN apk add --no-cache git build-base linux-headers

WORKDIR /app

# download and cache go mod
COPY ./go.* ./
RUN go env -w GO111MODULE=on && go mod download

COPY . .

RUN make build

# build pundix
FROM alpine:3.16

WORKDIR root

COPY --from=builder /app/build/bin/pundixd /usr/bin/pundixd

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp

VOLUME ["/root"]

ENTRYPOINT ["pundixd"]
