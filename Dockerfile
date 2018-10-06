# STEP 1: Used to get SSL root certificates
FROM alpine:3.8 as builder

# Install SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

# STEP 2: Create the abylebotter image
FROM scratch

# Copy SSL root certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY bin/abylebotter /usr/bin/

CMD ["/usr/bin/abylebotter", "-c", "/app/config/config.toml"]

LABEL org.label-schema.vendor="Abyle.org" \
      org.label-schema.url="https://github.com/torlenor/AbyleBotter" \
      org.label-schema.name="Abylebotter" \
      org.label-schema.description="An extensible chat bot for Discord written in GO"

