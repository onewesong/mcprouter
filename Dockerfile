FROM alpine

WORKDIR /data

COPY ./.env.example.toml ./.env.toml
COPY ./mcprouter ./

EXPOSE 8025

RUN chmod +x mcprouter

CMD ["./mcprouter", "server"]