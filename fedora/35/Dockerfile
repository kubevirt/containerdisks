FROM alpine as builder
ENV VERSION=35
RUN apk add wget jq
RUN wget https://getfedora.org/releases.json
RUN wget -O disk.qcow2 $(cat releases.json | jq -r --arg VERSION $VERSION '.[] | select(.link|test(".*qcow2")) | select(.variant=="Cloud" and .arch=="x86_64" and .version==$VERSION).link')
RUN echo "$(cat releases.json | jq -r --arg VERSION $VERSION '.[] | select(.link|test(".*qcow2")) | select(.variant=="Cloud" and .arch=="x86_64" and .version==$VERSION).sha256')  disk.qcow2" | sha256sum -c

FROM scratch                                                                                                    
COPY --chown=107:107 --from=builder /disk.qcow2 /disk/disk.qcow2
