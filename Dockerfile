FROM golang:1.10-alpine AS builder


ENV APP_USER ngtd-user
ENV APP_NAME ngtd

ENV NGT_VERSION 1.3.1

RUN set -eux \
    && apk --no-cache add libstdc++ ca-certificates \
    && apk --no-cache add --virtual build-dependencies cmake g++ make unzip curl upx git

RUN adduser -D -g '' ${APP_USER}

RUN curl -sSL "https://github.com/yahoojapan/NGT/archive/v${NGT_VERSION}.zip" -o NGT.zip \
    && unzip NGT.zip \
    && cd NGT-${NGT_VERSION} \
    && cmake . \
    && make -j \
    && make install \
    && cd .. \
    && rm -rf NGT.zip NGT-${NGT_VERSION}

RUN go get -v -u github.com/yahoojapan/ngtd \
    && go get -v -u github.com/mattn/go-sqlite3 \
    && go get -v -u gopkg.in/urfave/cli.v1

WORKDIR ${GOPATH}/src/github.com/yahoojapan/ngtd/cmd/ngtd

RUN CGO_ENABLED=1 \
    CGO_CXXFLAGS="-g -Ofast -march=native" \
    CGO_FFLAGS="-g -Ofast -march=native" \
    CGO_LDFLAGS="-g -Ofast -march=native" \
    GOOS=$(go env GOOS) \
    GOARCH=$(go env GOARCH) \
    go build --ldflags '-s -w -linkmode "external" -extldflags "-static -fPIC -m64 -pthread -std=c++17 -lstdc++"' -a -tags "cgo netgo" -installsuffix "cgo netgo" -o ${APP_NAME} \
    && upx -9 -o /usr/bin/${APP_NAME} ${APP_NAME}

RUN apk del build-dependencies --purge \
    && rm -rf ${GOPATH}

# Start From Scratch For Running Environment
FROM alpine:latest
LABEL maintainer "kpango <i.can.feel.gravity@gmail.com>"

ENV APP_USER ngtd-user
ENV APP_NAME ngtd

# Copy certificates for SSL/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy permissions
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}

EXPOSE 8080

USER ${APP_USER}

ENTRYPOINT ["/go/bin/ngtd"]
CMD ["grpc","-i","/tmp/index","-d","128","-t","bolt","-p","/tmp/kvs","-P","8080"]
