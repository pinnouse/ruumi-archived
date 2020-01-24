FROM golang:alpine

EXPOSE 8000

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN apk add git

RUN go get -d -v  .
RUN go build -o ruumi .

CMD ["/app/ruumi"]