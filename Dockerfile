FROM google/golang
MAINTAINER Shaalx Shi "60026668.m@daocloud.io"

# Build app
WORKDIR /gopath/app
ENV GOPATH /gopath/app
ADD . /gopath/app/

RUN go get -u github.com/everfore/oauth
RUN go install github.com/everfore/oauth

EXPOSE 80

CMD ["/gopath/app/bin/oauth"]