FROM alpine:latest

RUN mkdir -p /app
RUN mkdir -p /app/data
RUN mkdir -p /app/config

WORKDIR /app
COPY satellite /app/
COPY satellite.settings.json /app/config/

RUN chmod +x /app/satellite

ARG LOGGER_MODE=info
ARG LOGGER_SYNC_ACTIVE=true
ARG SATELLITE_AUTH_TOKEN=$(openssl rand -hex 16)

EXPOSE 9623

# Run the image
CMD ["/app/satellite"]