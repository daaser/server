FROM golang:1.14.3

WORKDIR /go/src/app

COPY . /go/src/app

ARG DEV
# this is only b/c of the stupid Corporate network
# when on a local box pass `--build-arg DEV=true` to docker build
RUN if [ "x${DEV}" = "xtrue" ]; then \
      cp /go/src/app/config/cacert.pem /etc/ssl/certs/; \
    fi && \
    make && \
    mv ./build/app /go/bin/app

ENTRYPOINT ["/go/bin/app"]

