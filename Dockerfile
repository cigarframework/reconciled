FROM k8s.gcr.io/debian-base:v1.0.0

ADD bin/rd-server /app/rd-server
RUN chmod a+x /app/rd-server
CMD ["/app/rd-server", "--config", "/app/config.yaml"]
