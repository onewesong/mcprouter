FROM alpine

WORKDIR /data

COPY ./.env.example.toml ./.env.toml
COPY ./mcprouter ./

EXPOSE 8025 8027

RUN chmod +x mcprouter

ENTRYPOINT ["./mcprouter"]

# proxy: 8025
# api: 8027
CMD ["proxy"] 