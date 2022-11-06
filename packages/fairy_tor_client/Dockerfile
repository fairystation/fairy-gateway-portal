# The base go-image
FROM golang:1.18.8-alpine3.16
 
# Create a directory for the app
RUN mkdir /app
 
# Set working directory
WORKDIR /app
 
# Run command as described:
# go build will build an executable file named server in the current directory
RUN apk add git; \
	go get -u gorm.io/gorm; \
	go get -u gorm.io/driver/sqlite; \
	go get -u github.com/gorilla/mux; \
	go get -u github.com/go-sql-driver/mysql; \
	go get github.com/canthefason/go-watcher; \
	go get -u github.com/andrewsjg/go-enigma; \
	go get -u github.com/auyer/steganography; \
	go install github.com/canthefason/go-watcher/cmd/watcher;

RUN apk add linux-headers gcc libc-dev openssl-dev libevent-dev zlib-dev;

RUN  CGO_CFLAGS="-w" go get -u -v -x github.com/cretz/bine/tor;\
     CGO_CFLAGS="-w" go get -u berty.tech/go-libtor;
#     CGO_CFLAGS="-w" go get -u github.com/ipsn/go-libtor;

RUN go get -u github.com/reiver/go-oi; \
    go get -u github.com/reiver/go-telnet;

RUN go get -u github.com/sadlil/go-trigger;

RUN go clean -modcache;

 
# Run the server executable
CMD [ "watcher" ]