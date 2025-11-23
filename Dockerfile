FROM ubuntu:latest
LABEL authors="varya"

ENTRYPOINT ["top", "-b"]