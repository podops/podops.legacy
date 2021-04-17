FROM registry.access.redhat.com/ubi8-minimal

ENV BUILD_NAME svc
ENV PROJECT_ID podops
ENV GOOGLE_APPLICATION_CREDENTIALS /config/credentials.json

EXPOSE 8080

# Copy extra files to the image
COPY ./run /root/usr/bin/

VOLUME /config
VOLUME /data

# copy the binary
COPY $BUILD_NAME /usr/bin/$BUILD_NAME

# make it executable
RUN chmod +x /usr/bin/$BUILD_NAME

CMD ["/root/usr/bin/run"]
