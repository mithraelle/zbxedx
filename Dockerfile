FROM golang:1.26-trixie
WORKDIR /build

ARG PLUGIN_IMPORT_SRC=build/zbxedx_plugins_linux.go
ARG PLUGIN_IMPORT_DST=/build/zabbix/src/go/plugins/zbxedx_plugins_linux.go

ARG CONFIG_IMPORT_SRC=conf/

ARG ZBX_VERSION=
ARG FORCE_BOOTSTRAP=0

RUN apt-get update && apt-get install -y --no-install-recommends \
    autoconf automake libtool pkg-config gettext \
    make gcc g++ \
    libpcre2-dev \
    zlib1g-dev \
    git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY ./ /build/

# Bootstrap Zabbix if the source tree is missing from the build context, or on FORCE_BOOTSTRAP=1
RUN if [ "${FORCE_BOOTSTRAP}" = "1" ] || [ ! -f zabbix/src/go/go.mod ]; then \
        bash ./bootstrap.sh --force ${ZBX_VERSION}; \
    else \
        echo "==> Using Zabbix source tree from build context"; \
    fi

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