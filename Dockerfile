FROM golang:1.17-alpine3.13

ENV HOME /go/src/readwrite
ENV PORT 8080
WORKDIR ${HOME}

COPY . ${HOME}/

ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

RUN go get -u github.com/go-redis/redis
RUN apk add --no-cache make bash
RUN GOOS=linux GOARCH=amd64 go build -o readwrite .
EXPOSE 8080
CMD [ "./readwrite" ]




