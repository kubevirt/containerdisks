FROM alpine as builder
ENV VERSION=8
ENV BASE_URL=https://cloud.centos.org/centos/$VERSION/x86_64/images/
RUN apk add wget
RUN wget -O - $BASE_URL | grep "GenericCloud.*qcow2" | sed -e 's/.*>\(.*qcow2\)<.*/\1/' | sort -r | head -n 1 > OS_IMAGE
RUN wget -q "${BASE_URL}$(cat OS_IMAGE)"
RUN wget -O - "${BASE_URL}CHECKSUM" | grep $(cat OS_IMAGE) | sed -ne 's/^SHA256 (\(.*\)) = \(.*\)/\2  \1/p' | sha256sum -c
RUN mv $(cat OS_IMAGE) disk.qcow2

FROM scratch
COPY --chown=107:107 --from=builder /disk.qcow2 /disk/disk.qcow2
