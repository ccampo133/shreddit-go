FROM gcr.io/distroless/static-debian12

COPY shreddit /

ENTRYPOINT ["/shreddit"]
