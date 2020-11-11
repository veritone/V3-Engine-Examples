# This uses a multi-stage build to get the engine toolkit
# see: https://docs.docker.com/develop/develop-images/multistage-build/)

# include the engine toolkit base image
FROM veritone/aiware-engine-toolkit as engine-toolkit

FROM alpine:latest

# get the engine-toolkit from the veritone base
COPY --from=engine-toolkit /opt/aiware/engine /opt/aiware/engine

# install the requirements for the engine-toolkit to run
RUN apk update \
    && apk --no-cache add ca-certificates libc6-compat

# add our engine and manifest
ADD dist/exif-extraction-engine /app/exif-extraction-engine
ADD manifest.json /var/manifest.json

# set up the environment variables
ENV VERITONE_WEBHOOK_READY="http://0.0.0.0:8080/ready"
ENV VERITONE_WEBHOOK_PROCESS="http://0.0.0.0:8080/process"

EXPOSE 8080

ENTRYPOINT [ "/opt/aiware/engine", "/app/exif-extraction-engine" ]
