FROM golang:alpine

EXPOSE 8000

RUN mkdir /app
ADD . /app
WORKDIR /app

ENV S3_ENDPOINT url_of_endpoint
ENV S3_BUCKET name_of_bucket
ENV AWS_ACCESS_KEY_ID your_access_key
ENV AWS_SECRET_ACCESS_KEY your_secret_key

RUN apk add git

RUN go get -d -v  .
RUN go build -o ruumi .

CMD ["/app/ruumi"]