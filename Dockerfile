FROM alpine:3.7
MAINTAINER chenyong scchenyong@189.cn
ENV REFRESHED_AT 2018-07-07
COPY  ./tcpecho   /usr/bin
RUN   chmod +x    /usr/bin/tcpecho
ENTRYPOINT  ["/usr/bin/tcpecho"]

