# platforms=linux/amd64
# image=ghcr.io/cidverse/rewrite-app

# runtime image
#
FROM quay.io/cidverse/base-coretto-21:21.0.2.13.1

ADD rootfs /
ADD .dist/github-com-cidverse-rewrite-app/binary/linux_amd64 /usr/local/bin/rewrite-app

RUN chmod +x /usr/local/bin/*

CMD ["bash"]
