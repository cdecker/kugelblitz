FROM ubuntu:16.04
MAINTAINER Christian Decker <decker@blockstream.com>
RUN apt-get update -q && apt-get install -y supervisor git software-properties-common && \
    add-apt-repository ppa:bitcoin/bitcoin && \
    apt-get update -q && \
    apt-get install -y \
    	    autoconf \
	    libtool \
	    libprotobuf-c-dev \
	    libsodium-dev \
	    libsqlite3-dev \
	    libgmp-dev \
	    build-essential \
	    libsqlite3-dev \
	    bitcoind \
	    protobuf-c-compiler \
	    g++ \
	    curl \
	    gcc \
	    libc6-dev \
	    make \
	    pkg-config && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get install -y build-essential
RUN git clone https://github.com/luke-jr/libbase58.git /opt/libbase58 && cd /opt/libbase58 && ./autogen.sh && ./configure && make && make install
RUN useradd lightning; mkdir /lightning /bitcoin; chown lightning:users /lightning /bitcoin
RUN git clone https://github.com/ElementsProject/lightning.git /opt/lightning; cd /opt/lightning; make

ADD supervisor.conf /etc/supervisor/conf.d/
RUN echo "testnet=1\nrpcuser=rpcuser\nrpcpassword=rpcpass" > /bitcoin/bitcoin.conf

RUN sed -i 's/^\(\[supervisord\]\)$/\1\nnodaemon=true/' /etc/supervisor/supervisord.conf
CMD ["supervisord", "-c", "/etc/supervisor/supervisord.conf"]


RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.7.3
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 508028aac0654e993564b6e2014bf2d4a9751e3b286661b0b0040046cf18028e

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN go get github.com/cdecker/kugelblitz/...

VOLUME ['/bitcoin', '/lightning']
EXPOSE 19735
