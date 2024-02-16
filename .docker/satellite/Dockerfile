FROM alpine:latest

RUN mkdir -p /app
RUN mkdir -p /app/data

WORKDIR /app/data
COPY /dist/satellite /app

CMD ["/app/satellite"]