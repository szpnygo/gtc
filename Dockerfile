FROM alpine

WORKDIR /app
COPY gtc /app/gtc

CMD ["/app/gtc"]
