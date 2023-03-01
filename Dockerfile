FROM golang:1.20.1-alpine3.16 as builder

RUN apk update && apk upgrade && apk add --no-cache curl git

RUN curl -s https://raw.githubusercontent.com/eficode/wait-for/master/wait-for -o /usr/bin/wait-for
RUN chmod +x /usr/bin/wait-for

RUN go install github.com/prep/pubsubc@latest

###############################################################################

FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:420.0.0-alpine

COPY --from=builder /usr/bin/wait-for /usr/bin
COPY --from=builder /go/bin/pubsubc   /usr/bin
COPY                run.sh            /run.sh

RUN apk --update add openjdk8-jre netcat-openbsd && gcloud components install beta pubsub-emulator

EXPOSE 8681

CMD /run.sh
