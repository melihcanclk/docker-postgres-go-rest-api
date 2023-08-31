FROM golang:1.21.0

WORKDIR /usr/src/app

# reason why we do not include this by getting is that we want to call this
# using command line command `air` instead of go run
RUN go install github.com/cosmtrek/air@latest 

COPY . .

RUN go mod tidy