FROM golang:latest 
# Add Maintainer Info
LABEL maintainer="Natnael Awel <natnael.awel@gmail.com>"

ADD . /go/src/github.com/tenahubclientdocker
WORKDIR /go/src/github.com/tenahubclientdocker

RUN go get ./... 
#RUN go get -d -v ./...
RUN go install /go/src/github.com/tenahubclientdocker

ENTRYPOINT /go/bin/tenahubclientdocker

#CMD ["tenahubclientdocker"]
EXPOSE 8282

