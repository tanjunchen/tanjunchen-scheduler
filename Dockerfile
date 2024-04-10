FROM debian:stretch-slim

WORKDIR /

COPY _output/bin/tanjunchen-scheduler /usr/local/bin

CMD ["tanjunchen-scheduler"]