FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        clamav-daemon \
        ca-certificates \
        gettext-base && rm -rf /var/lib/apt/lists/*
RUN freshclam --user clamav
EXPOSE 3310
RUN rm -rf /var/log/clamav/* && \
    ln -s /dev/stdout /var/log/clamav/clamav.log && \
    rm -f /etc/clamav/clamd.conf && \
    mkdir /var/run/clamav && \
    chown clamav:clamav /var/run/clamav
COPY conf.d/clamav.conf.template /etc/clamav/clamd.conf.template
ENV CLAMD_LISTENING_PORT=3310
CMD envsubst < /etc/clamav/clamd.conf.template > /etc/clamav/clamd.conf && clamd