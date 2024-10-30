FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN ln -fs /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apt update && \
    apt install -y clamav-freshclam curl unzip && \
    dpkg-reconfigure --frontend noninteractive tzdata && \
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install && \
    mkdir -p /tmp/clamdb && \
    chown -R 101:102 /tmp/clamdb

COPY freshclam.conf /etc/clamav/freshclam.conf
COPY downloader.sh /usr/local/bin/downloader.sh
RUN chmod +x /usr/local/bin/downloader.sh

CMD ["/usr/local/bin/downloader.sh"]