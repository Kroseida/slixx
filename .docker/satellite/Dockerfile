FROM ubuntu:latest

RUN mkdir -p /app
RUN mkdir -p /app/data
RUN mkdir -p /app/config
RUN mkdir -p /app/log

WORKDIR /app
COPY satellite /app/
COPY satellite.settings.json /app/config/

ENV LOGGER_MODE=info
ENV LOGGER_SYNC_ACTIVE=true
ENV SATELLITE_AUTH_TOKEN=

EXPOSE 9623

# Run the image
ENTRYPOINT ["/app/satellite"]