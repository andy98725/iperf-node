
FROM golang:1.21.6 AS build
WORKDIR /

COPY src/go.mod ./go.mod
COPY src/go.sum ./go.sum
RUN go mod download
COPY src/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-server

FROM alpine

### ==================================
ENV ENDPOINT=http://host.docker.internal:3001
ENV ENDPOINT_KEY="<enter key here>"
ENV HASH="<enter hash here>"
ENV ID="<enter id here>"

ENV IPERF_PORT=5001
ENV PORT=8080
EXPOSE 5001
EXPOSE 8080

LABEL AUTHOR=andyhudson725@gmail.com
LABEL VERSION=0.1
### ==================================

RUN apk update && \
    apk add curl && \
    apk add tar && \
    # To compile iperf
    apk add build-base


# Install Iperf
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

CMD /go-server
