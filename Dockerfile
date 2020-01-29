FROM golang:alpine

EXPOSE 8000

RUN mkdir /app
ADD . /app
WORKDIR /app

ENV S3_ENDPOINT url_of_endpoint
ENV S3_BUCKET name_of_bucket

RUN apk add git

RUN go get -d -v  .
RUN go build -o ruumi .

CMD ["/app/ruumi"]