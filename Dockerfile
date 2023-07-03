FROM scratch
LABEL org.opencontainers.image.authors="Cybozu, Inc." \
      org.opencontainers.image.title="pod-security-admission" \
      org.opencontainers.image.source="https://github.com/cybozu-go/pod-security-admission"
WORKDIR /
COPY pod-security-admission /
USER 10000:10000

ENTRYPOINT ["/pod-security-admission"]
