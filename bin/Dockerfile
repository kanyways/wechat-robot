FROM alpine:latest
RUN apk --no-cache add tzdata ca-certificates && cp -r -f /usr/share/zoneinfo/Hongkong /etc/localtime
WORKDIR /root
COPY wechat_robot /root/
EXPOSE 8899
ENTRYPOINT /root/wechat_robot -c config.docker.yaml