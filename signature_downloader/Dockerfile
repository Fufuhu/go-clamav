FROM debian:trixie-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates clamav-freshclam curl tzdata unzip && \
    ln -fs /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    dpkg-reconfigure --frontend noninteractive tzdata && \
    curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o awscliv2.zip && \
    unzip awscliv2.zip && \
    ./aws/install && \
    rm -rf awscliv2.zip aws && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /tmp/clamdb && \
    chown -R clamav:clamav /tmp/clamdb

COPY freshclam.conf /etc/clamav/freshclam.conf
COPY downloader.sh /usr/local/bin/downloader.sh
RUN chmod +x /usr/local/bin/downloader.sh

CMD ["/usr/local/bin/downloader.sh"]
