FROM golang:alpine AS build-env
WORKDIR /app
ADD . /app
RUN cd /app && env go build -o kubewise
FROM alpine
RUN apk update && \
   apk add ca-certificates && \
   update-ca-certificates && \
   rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/kubewise /app
ENTRYPOINT ./kubewise
