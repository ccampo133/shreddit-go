FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY shreddit /bin/shreddit

ENTRYPOINT ["/bin/shreddit"]
