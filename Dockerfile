FROM ubuntu:latest AS builder

ENV APP_NAME ngtd

ENV NGT_VERSION 1.7.3

ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG ja_JP.UTF-8
ENV GOVERSION 1.12.5
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

WORKDIR ${GOPATH}/src/github.com/yahoojapan/ngtd
COPY . .

RUN make deps

WORKDIR ${GOPATH}/src/github.com/yahoojapan/ngtd/cmd/ngtd
RUN CGO_ENABLED=1 \
    CGO_CXXFLAGS="-g -Ofast -march=native" \
    CGO_FFLAGS="-g -Ofast -march=native" \
    CGO_LDFLAGS="-g -Ofast -march=native" \
    GOOS=$(go env GOOS) \
    GOARCH=$(go env GOARCH) \
    GO111MODULE=on \
    go build --ldflags '-s -w -linkmode "external" -extldflags "-static -fPIC -m64 -pthread -fopenmp -std=c++17 -lstdc++ -lm"' -a -tags "cgo netgo" -installsuffix "cgo netgo" -o ${APP_NAME} \
    # && upx -9 -o /usr/bin/${APP_NAME} ${APP_NAME}
    && upx --best --ultra-brute -o /usr/bin/${APP_NAME} ${APP_NAME}

# Start From Scratch For Running Environment
FROM scratch
LABEL maintainer "kpango <i.can.feel.gravity@gmail.com>, Kosuke Morimoto <kou.morimoto@gmail.com>"

ENV APP_NAME ngtd

# Copy certificates for SSL/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy permissions
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /usr/bin/${APP_NAME} /${APP_NAME}

EXPOSE 8200

ENTRYPOINT ["/ngtd"]
