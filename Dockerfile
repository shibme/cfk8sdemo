FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
WORKDIR /
COPY ./dist/app_linux_${TARGETARCH}*/app .
USER 65532:65532
ENTRYPOINT ["/app"]