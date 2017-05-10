FROM busybox

ENV APP_HOME=/app
ENV APP_DATA=/app/data
ENV APP_LOGS=/app/logs

RUN mkdir -p ${APP_DATA}
RUN mkdir -p ${APP_LOGS}

COPY ./data/ ${APP_DATA}/
COPY server ${APP_HOME}/
COPY start.sh ${APP_HOME}/

EXPOSE 8087

WORKDIR ${APP_HOME}

RUN chmod +x server
RUN chmod +x start.sh

CMD /app/start.sh