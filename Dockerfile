FROM scratch

COPY bin/abylebotter  /usr/bin/
CMD ["/usr/bin/abylebotter"]

LABEL org.label-schema.vendor="Abyle.org" \
      org.label-schema.url="https://github.com/torlenor/AbyleBotter" \
      org.label-schema.name="Abylebotter" \
      org.label-schema.description="An extensible chat bot for Discord written in GO" 