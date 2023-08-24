FROM        alpine

COPY ./bin/gofiber.linux /usr/local/main/main
COPY ./.env /usr/local/main/.env
CMD mkdir -p /usr/local/main/webdata
#ADD ./assets /usr/local/main/assets
#ADD ./views /usr/local/main/views

WORKDIR /usr/local/main
CMD ./main