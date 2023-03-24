FROM golang:latest
ENV LANGUAGE="en"

WORKDIR /app
COPY ./ /app

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o telegramtz ./main.go

CMD ["./telegramtz"]