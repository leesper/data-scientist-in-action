FROM leesper/ubuntu-java

ENV ZOOKEEPER_VERSION=3.4.10
ENV ZOOKEEPER_URL=http://mirrors.shu.edu.cn/apache/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/zookeeper-${ZOOKEEPER_VERSION}.tar.gz

RUN set -x \
    && curl -fSL "${ZOOKEEPER_URL}" -o /tmp/zookeeper.tar.gz \
    && tar -xvf /tmp/zookeeper.tar.gz -C /usr/local \
    && mv /usr/local/zookeeper-${ZOOKEEPER_VERSION} /usr/local/zookeeper \
    && rm -rf /tmp/zookeeper.tar.gz* \
    && mkdir -p /opt/data \
    && mkdir -p /opt/log

ENV ZOO_HOME=/usr/local/zookeeper
ENV PATH=$PATH:${ZOO_HOME}/bin
ENV TZ="Asia/Shanghai"
EXPOSE 2181 2888 3888
CMD ["zkServer.sh", "start-foreground"]