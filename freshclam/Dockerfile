FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        clamav-freshclam \
        ca-certificates \
        gettext-base && rm -rf /var/lib/apt/lists/*
RUN ln -s /dev/stdout /var/log/clamav/clamav.log

COPY conf.d/freshclam.conf.template /etc/clamav/freshclam.conf.template
COPY src/run.sh /usr/local/bin/run.sh
RUN chmod +x /usr/local/bin/run.sh && \
    chmod 644 /etc/clamav/freshclam.conf && \
    chown 100:102 /var/lib/clamav 

# read-only root filesystem 対応: 設定生成と freshclam 実行は run.sh 内で行う
CMD ["run.sh"]
