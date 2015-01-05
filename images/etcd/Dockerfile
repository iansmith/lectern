FROM phusion/baseimage:0.9.15
MAINTAINER eng@igneous.io

WORKDIR /tmp
RUN curl -L  https://github.com/coreos/etcd/releases/download/v0.4.6/etcd-v0.4.6-linux-amd64.tar.gz -o etcd-v0.4.6-linux-amd64.tar.gz
RUN tar xzvf etcd-v0.4.6-linux-amd64.tar.gz

RUN mkdir -p /etc/service/etcd
RUN cp /tmp/etcd-v0.4.6-linux-amd64/etcd /etc/service/etcd/etcd
RUN chmod +x /etc/service/etcd/etcd

ADD run-etcd /etc/service/etcd/run
RUN chmod +x /etc/service/etcd/run

RUN cp /tmp/etcd-v0.4.6-linux-amd64/etcdctl /bin/etcdctl
RUN chmod +x /etc/service/etcd/run
RUN chmod +x /bin/etcdctl

ADD etcd_config /etc/my_init.d/etcd_config
RUN chmod +x /etc/my_init.d/etcd_config

EXPOSE 4001
EXPOSE 7001
ENTRYPOINT ["/sbin/my_init"]

