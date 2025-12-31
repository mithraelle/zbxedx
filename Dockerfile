FROM golang:1.24-bookworm
WORKDIR /build

ARG PLUGIN_IMPORT_SRC=build/zbxedx_plugins_linux.go
ARG PLUGIN_IMPORT_DST=/build/zabbix/src/go/plugins/zbxedx_plugins_linux.go

ARG CONFIG_IMPORT_SRC=conf/

RUN apt-get update && apt-get install -y --no-install-recommends \
    autoconf automake libtool pkg-config gettext \
    make gcc g++ \
    libpcre2-dev \
    zlib1g-dev \
    git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY ./ /build/
COPY ${PLUGIN_IMPORT_SRC} ${PLUGIN_IMPORT_DST}

RUN mkdir -p /build/zabbix/src/go/zbxedx \
    && cp -r /build/edxagent /build/zabbix/src/go/zbxedx/edxagent \
    && cp -r /build/plugins /build/zabbix/src/go/zbxedx/plugins

WORKDIR /build/zabbix
RUN ./bootstrap.sh
RUN ./configure --enable-agent2
RUN make
RUN make install

RUN rm -rf /build

COPY ${CONFIG_IMPORT_SRC} /usr/local/etc/

CMD ["zabbix_agent2", "-f", "-c", "/usr/local/etc/zabbix_agent2.conf"]