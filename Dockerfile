FROM scratch
RUN apk update && \
   apk add ca-certificates && \
   update-ca-certificates && \
   rm -rf /var/cache/apk/*
COPY kubewise /
ENTRYPOINT ["/kubewise"]

