FROM golang:1.21.4-alpine as builder

RUN apk update && apk upgrade && apk add --no-cache curl git

RUN curl -s https://raw.githubusercontent.com/eficode/wait-for/master/wait-for -o /usr/bin/wait-for
RUN chmod +x /usr/bin/wait-for

RUN go install github.com/prep/pubsubc@latest

###############################################################################

FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:454.0.0-alpine

COPY --from=builder /usr/bin/wait-for /usr/bin
COPY --from=builder /go/bin/pubsubc   /usr/bin
COPY                run.sh            /run.sh

ENV LD_PRELOAD=/lib/libgcompat.so.0
RUN apk --update add openjdk8-jre netcat-openbsd gcompat && gcloud components install beta pubsub-emulator

EXPOSE 8681

CMD /run.sh
