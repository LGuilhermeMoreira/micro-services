FROM ubuntu:latest
LABEL authors="guigui"

ENTRYPOINT ["top", "-b"]