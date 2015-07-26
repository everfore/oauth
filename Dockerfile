FROM google/golang
MAINTAINER Shaalx Shi "60026668.m@daocloud.io"

# Build app
WORKDIR /gopath/app
ENV GOPATH /gopath/app
ADD . /gopath/app/

RUN go get -u github.com/shaalx/echo
RUN go install github.com/shaalx/echo
#RUN go build -o echo
RUN ls
EXPOSE 80
CMD ["/gopath/app/bin/echo"]