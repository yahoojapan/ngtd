FROM ubuntu:16.04 AS builder


ENV APP_NAME ngtd

ENV NGT_VERSION 1.3.2

ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG ja_JP.UTF-8
ENV GOVERSION 1.10.1
ENV GOROOT /opt/go
ENV GOPATH /go

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    cmake \
    upx \
    curl \
    unzip \
    git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN cd /opt && curl -sSL -O https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH

RUN curl -sSL "https://github.com/yahoojapan/NGT/archive/v${NGT_VERSION}.zip" -o NGT.zip \
    && unzip NGT.zip \
    && cd NGT-${NGT_VERSION} \
    && cmake . \
    && make -j \
    && make install \
    && cd .. \
    && rm -rf NGT.zip NGT-${NGT_VERSION}

RUN go get -v -u github.com/golang/dep/cmd/dep \
    && dep ensure

WORKDIR ${GOPATH}/src/github.com/yahoojapan/ngtd/cmd/ngtd

RUN CGO_ENABLED=1 \
    CGO_CXXFLAGS="-g -Ofast -march=native" \
    CGO_FFLAGS="-g -Ofast -march=native" \
    CGO_LDFLAGS="-g -Ofast -march=native" \
    GOOS=$(go env GOOS) \
    GOARCH=$(go env GOARCH) \
    go build --ldflags '-s -w -linkmode "external" -extldflags "-static -fPIC -m64 -lm -pthread -std=c++17 -lstdc++"' -a -tags "cgo netgo" -installsuffix "cgo netgo" -o ${APP_NAME} \
    && upx -9 -o /usr/bin/${APP_NAME} ${APP_NAME}

# Start From Scratch For Running Environment
FROM scratch
LABEL maintainer "kpango <i.can.feel.gravity@gmail.com>, Kosuke Morimoto <kou.morimoto@gmail.com>"

ENV APP_USER ngtd-user
ENV APP_NAME ngtd

# Copy certificates for SSL/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy permissions
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /usr/bin/${APP_NAME} /${APP_NAME}

EXPOSE 8200

ENTRYPOINT ["/ngtd"]
