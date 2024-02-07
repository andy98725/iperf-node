
FROM golang:1.21.6 AS build
WORKDIR /

COPY src/go.mod ./go.mod
COPY src/go.sum ./go.sum
RUN go mod download
COPY src/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-server

FROM debian:stable-slim

RUN apt update && \
    apt install curl -y && \
    apt install tar -y && \
    # To compile iperf
    apt install gcc -y && \
    apt install g++ -y && \
    apt install make -y

# Install Iperf
RUN curl -L https://sourceforge.net/projects/iperf2/files/iperf-2.1.9.tar.gz/download \
    | tar -zxv -C /tmp
RUN cd /tmp/iperf-2.1.9 && \
    ./configure && \
    make install

# Cleanup
RUN rm -rf /tmp/iperf-2.1.9 && \
    apt remove curl -y && \
    apt remove gcc -y && \
    apt remove g++ -y && \
    apt remove make -y


# Setup golang server
WORKDIR /
COPY --from=build /go-server /go-server

RUN adduser nonroot 
USER nonroot:nonroot

### ==================================
ENV ENDPOINT=https://iperf-benchmark-stg-78bf74de879a.herokuapp.com
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

CMD /go-server
