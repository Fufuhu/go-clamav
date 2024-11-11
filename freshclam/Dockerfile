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

CMD envsubst < /etc/clamav/freshclam.conf.template > /etc/clamav/freshclam.conf && run.sh
