FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        clamav-daemon \
        ca-certificates \
        gettext-base && rm -rf /var/lib/apt/lists/*
EXPOSE 3310
RUN rm -rf /var/log/clamav/* && \
    ln -s /dev/stdout /var/log/clamav/clamav.log && \
    rm -f /etc/clamav/clamd.conf && \
    mkdir /var/run/clamav && \
    chown clamav:clamav /var/run/clamav

COPY conf.d/clamav.conf.template /etc/clamav/clamd.conf.template
COPY src/run.sh /usr/local/bin/run.sh
RUN chmod +x /usr/local/bin/run.sh

ENV CLAMD_LISTENING_PORT=3310
# read-only root filesystem 対応のため、設定ファイルの生成は run.sh 内で
# 書き込み可能な /run/clamav に対して行う
CMD ["run.sh"]