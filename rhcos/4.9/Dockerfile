FROM alpine as builder
RUN wget https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.9/latest/rhcos-openstack.x86_64.qcow2.gz
RUN wget -O - https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.9/latest/sha256sum.txt | grep rhcos-openstack.x86_64.qcow2.gz | sha256sum -c
RUN gzip -d rhcos-openstack.x86_64.qcow2.gz

FROM scratch                                                                                                    
COPY --chown=107:107 --from=builder /rhcos-openstack.x86_64.qcow2 /disk/rhcos-openstack.x86_64.qcow2
