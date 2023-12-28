FROM golang:1.21.5

WORKDIR github.com/jamesmukumu/diary2024

COPY  go.mod go.sum ./

RUN go mod download

COPY  . .

EXPOSE 6600

CMD [ "go","run","main.go" ]