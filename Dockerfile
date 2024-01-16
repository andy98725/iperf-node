
# FROM httpd:2.4
FROM alpine

LABEL AUTHOR=andyhudson725@gmail.com
LABEL VERSION=0.1

ENV SERVER_PORT=5001
EXPOSE 5001

RUN apk update && \
    apk add curl && \
    apk add tar && \
    # To compile iperf
    apk add build-base


# Install iPerf
RUN curl -L https://sourceforge.net/projects/iperf2/files/iperf-2.1.9.tar.gz/download \
        | tar -zxv -C /tmp
RUN cd /tmp/iperf-2.1.9 && \
    ./configure && \
    make install

# Cleanup
RUN rm -rf /tmp/iperf-2.1.9 && \
    # needed by iperf
    # apk del build-base && \ 
    apk del curl && \
    apk del tar

WORKDIR /root/

ENTRYPOINT [ "/usr/local/bin/iperf" ]
CMD [ "-s -p 5001" ]
