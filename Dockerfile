FROM golang:1.4
MAINTAINER Justin C. Miller <justin@devjustinian.com>

ADD . /go/src/github.com/justinian/slackdice

RUN cd /go/src/github.com/justinian/slackdice && go get && go install

CMD ["/go/bin/slackdice"]

EXPOSE 8000
