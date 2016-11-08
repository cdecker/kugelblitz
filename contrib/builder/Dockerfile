FROM ubuntu:16.04
MAINTAINER Christian Decker <decker.christian@gmail.com>

RUN echo "deb http://ppa.launchpad.net/bitcoin/bitcoin/ubuntu xenial main" >> /etc/apt/sources.list && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys C70EF1F0305A1ADB9986DBD8D46F45428842CE5E
RUN apt-get update -qq && \
    apt-get install -y --no-install-recommends \
	bitcoind \
	build-essential \
	autoconf \
	automake \
	eatmydata \
	git \
	net-tools \
	libtool \
	libprotobuf-c-dev \
	libgmp-dev \
	libsodium-dev \
	libsqlite3-dev \
	libbase58-dev \
	valgrind && \
	apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD build.sh /opt/
VOLUME /build
CMD /opt/build.sh
