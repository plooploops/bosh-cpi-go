ARG GO_VERSION=1.10.1
ARG ALPINE_VERSION=3.7

######################################

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as BUILDENV

RUN apk add --update curl git gcc && \
    rm -rf /var/cache/apk/* && \
    go get -u github.com/golang/dep/cmd/dep

######################################

FROM BUILDENV as BUILDER

WORKDIR /go/src/github.com/plooploops/bosh-cpi-go
COPY ./ /go/src/github.com/plooploops/bosh-cpi-go

RUN dep ensure
# RUN go build -a -tags netgo -installsuffix netgo -ldflags '-w' docs/kubernetes-cpi.go
RUN go build docs/kubernetes-cpi.go



######################################

FROM apline:${ALPINE_VERSION}

#RUN apt-get -y update

COPY --from=BUILDER /go/src/github.com/chgeuer/GoSAS/CGI/server /usr/local/apache2/cgi-bin/server
