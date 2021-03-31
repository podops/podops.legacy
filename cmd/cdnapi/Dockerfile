FROM registry.access.redhat.com/ubi8-minimal

ENV BUILD_NAME cdnapi
#ENV CADDYFILE /config/caddy/Caddyfile

EXPOSE 8080

# See https://caddyserver.com/docs/conventions#file-locations for details
#ENV XDG_CONFIG_HOME=/config
#ENV XDG_DATA_HOME=/data

#RUN set -eux; \
#	mkdir -p \
#		/config/caddy \
#		/data/caddy \
#		/etc/caddy

# Copy extra files to the image
COPY ./run /root/usr/bin/

#VOLUME /config
#VOLUME /data

# copy the binary
COPY $BUILD_NAME /usr/bin/$BUILD_NAME

# make it executable
RUN chmod +x /usr/bin/$BUILD_NAME

CMD ["/root/usr/bin/run"]
