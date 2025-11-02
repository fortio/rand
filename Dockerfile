FROM scratch
COPY rand /usr/bin/rand
ENV HOME=/home/user
ENTRYPOINT ["/usr/bin/rand"]
