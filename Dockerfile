FROM alpine:latest

WORKDIR root

COPY ./build/bin/pundixd /usr/bin/pundixd

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp

VOLUME ["/root"]

ENTRYPOINT ["pundixd"]
