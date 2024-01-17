FROM golang:1.21.4

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x waitForPostgres.sh

# build go app
RUN go mod download
RUN go build -o apigame ./cmd/main.go

CMD [ "./apigame" ]