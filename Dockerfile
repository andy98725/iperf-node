
FROM golang:1.21.6 AS build
WORKDIR /

COPY src/ ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-server

FROM alpine

LABEL AUTHOR=andyhudson725@gmail.com
LABEL VERSION=0.1

ENV IPERF_PORT=5001
EXPOSE 5001
ENV PORT=8080
EXPOSE 8080

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


# Setup golang server
WORKDIR /
COPY --from=build /go-server /go-server

RUN adduser -D nonroot 
USER nonroot:nonroot

# ENTRYPOINT [ "/usr/local/bin/iperf" ]
# CMD [ "-s -p 5001" ]
CMD /go-server
