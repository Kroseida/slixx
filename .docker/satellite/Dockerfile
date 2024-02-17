FROM alpine:latest

RUN mkdir -p /app
RUN mkdir -p /app/data

WORKDIR /app
COPY /dist/satellite /app

ARG LOGGER_MODE=info
ARG LOGGER_SYNC_ACTIVE=true
ARG SATELLITE_AUTH_TOKEN=$(openssl rand -hex 16)

CMD ["satellite"]