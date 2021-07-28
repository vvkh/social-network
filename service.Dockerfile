FROM golang:1.16

COPY . /app
WORKDIR /app

RUN make build
CMD make up
