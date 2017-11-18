FROM alpine:3.4

RUN apk -U add ca-certificates

EXPOSE 8080

ADD news /bin/news
ADD config.yml /etc/news/config.yml

CMD ["news", "-config", "/etc/news/config.yml"]
