FROM alpine

# install kubectl
RUN apk add --no-cache curl && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/ && \
    rm -rf /var/cache/apk/*

WORKDIR /data

COPY ./.env.example.toml ./.env.toml
COPY ./mcprouter ./

EXPOSE 8025 8027

RUN chmod +x mcprouter

ENTRYPOINT ["./mcprouter"]

CMD ["proxy"]