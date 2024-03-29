FROM golang:1.20-alpine
ENV GOOS=linux
WORKDIR /app

COPY . .

WORKDIR  /app/

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

RUN chmod -R 777 .

EXPOSE 2222

ENTRYPOINT CompileDaemon -build="go build -o ironexalt" -command="./ironexalt"