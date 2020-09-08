FROM scratch

ADD --chown=1000:2000 https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY bin/feedy /

ENTRYPOINT ["/feedy"]
