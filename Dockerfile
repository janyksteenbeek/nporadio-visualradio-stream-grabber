FROM golang:latest

RUN apt-get update && apt-get install -y wget gnupg2 \
    && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list' \
    && apt-get update \
    && apt-get install -y google-chrome-stable

WORKDIR /app
COPY . /app

RUN go mod download

RUN go build -o nporadio-visualradio-stream-grabber cmd/grabber/main.go
RUN chmod +x nporadio-visualradio-stream-grabber

CMD ["/app/nporadio-visualradio-stream-grabber"]
