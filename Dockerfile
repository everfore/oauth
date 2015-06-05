FROM google/golang
MAINTAINER zhuulx "zhuulx@daocloud.io"

# Build app
WORKDIR /gopath/app
ENV GOPATH /gopath/app
ADD . /gopath/app/

RUN go get github.com/zhuulx/oauth
RUN go install github.com/zhuulx/oauth

EXPOSE 80
CMD ["/gopath/app/bin/oauth"]