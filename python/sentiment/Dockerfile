FROM veritone/aiware-engine-toolkit as vt-engine-toolkit
FROM ubuntu:18.04

RUN apt-get update \
    && apt-get install -y vim python3-pip python3-dev \
    && apt-get install -y build-essential cmake \
    && pip3 install Flask==1.1.2

# Install nltk and data packs
RUN pip3 install nltk \
    && python3 -m nltk.downloader vader_lexicon \
    && python3 -m nltk.downloader punkt

WORKDIR /app
COPY *.* ./

ENV PYTHONPATH="${PYTHONPATH}:/app"
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

ENV VERITONE_WEBHOOK_READY="http://0.0.0.0:8082/ready"
ENV VERITONE_WEBHOOK_PROCESS="http://0.0.0.0:8082/process"
ENV VERITONE_DISABLE_CHUNK_DOWNLOAD=true

ADD manifest.json /var/manifest.json

COPY --from=vt-engine-toolkit /opt/aiware/engine /opt/aiware/engine

EXPOSE 8082
ENTRYPOINT ["/opt/aiware/engine", "python3", "-m", "flask", "run", "--host=0.0.0.0", "--port=8082"]
