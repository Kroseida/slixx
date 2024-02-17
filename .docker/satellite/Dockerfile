FROM alpine:latest

RUN mkdir -p /app
RUN mkdir -p /app/data
RUN mkdir -p /app/config

WORKDIR /app
COPY ./dist/satellite /app/satellite
COPY ./.docker/satellite/satellite.settings.json /app/config/satellite.settings.json

ARG LOGGER_MODE=info
ARG LOGGER_SYNC_ACTIVE=true
ARG SATELLITE_AUTH_TOKEN=$(openssl rand -hex 16)

ENTRYPOINT ["satellite"]