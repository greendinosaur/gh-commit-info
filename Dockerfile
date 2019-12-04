FROM golang:alpine
RUN apk update && apk add git
EXPOSE 8080

# We copy everything in the root directory
# into our /app directory
ADD . /go/src/github.com/greendinosaur/gh-commit-info/
# We specify that we now wish to execute 
# any further commands inside our /app
# directory
WORKDIR /go/src/github.com/greendinosaur/gh-commit-info/src/api
# we run go build to compile the binary
# executable of our Go program
RUN go get -v
RUN go build -o main .

# Our start command which kicks off
# our newly created binary executable
CMD ["./main"]