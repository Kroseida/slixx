FROM alpine:latest

RUN mkdir -p /app
RUN mkdir -p /app/data

WORKDIR /app/data
COPY /dist/supervisor /app

CMD ["/app/supervisor"]