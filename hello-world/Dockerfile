FROM veritone/aiware-engine-toolkit as vt-engine-toolkit
FROM mhart/alpine-node:14.8

COPY . /app
COPY manifest.json /var/

WORKDIR /app

EXPOSE 8080

RUN apk update \
  && apk upgrade \
  && apk --no-cache add ca-certificates \
  && apk add --no-cache libc6-compat
RUN npm install

ENV VERITONE_WEBHOOK_READY="http://0.0.0.0:8080/readyz"
ENV VERITONE_WEBHOOK_PROCESS="http://0.0.0.0:8080/process"

COPY --from=vt-engine-toolkit /opt/aiware/engine /opt/aiware/engine

RUN ["chmod", "+x", "/app/index.js"]

ENTRYPOINT [ "/opt/aiware/engine", "node", "index.js" ]
