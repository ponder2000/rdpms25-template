FROM alpine:latest

RUN apk add --no-cache tzdata

ENV TZ=Asia/Calcutta
ENV GIN_MODE=release

WORKDIR /app

COPY bin/alpine/edge_playback .

CMD [ "./app" ]
