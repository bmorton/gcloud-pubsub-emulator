FROM golang:1.22 as builder

COPY . .
RUN go install .

###############################################################################

FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:454.0.0-alpine

ENV LD_PRELOAD=/lib/libgcompat.so.0
RUN apk --update add openjdk8-jre netcat-openbsd gcompat && gcloud components install beta pubsub-emulator

EXPOSE 8681

ENV PUBSUB_EMULATOR_HOST=localhost:8681

COPY --from=builder /go/bin/gcloud-pubsub-emulator /usr/bin

ENTRYPOINT ["gcloud-pubsub-emulator"]
