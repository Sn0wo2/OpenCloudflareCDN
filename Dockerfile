FROM alpine:latest

RUN mkdir -p /opt/OpenCloudflareCDN
WORKDIR /opt/OpenCloudflareCDN

COPY OpenCloudflareCDN /opt/OpenCloudflareCDN/OpenCloudflareCDN

RUN chmod +x /opt/OpenCloudflareCDN/OpenCloudflareCDN

ENTRYPOINT ["/opt/OpenCloudflareCDN/OpenCloudflareCDN"]
