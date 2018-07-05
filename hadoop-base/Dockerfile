FROM ubuntu:18.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends openjdk-8-jdk net-tools curl \
    && rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64

ENV HADOOP_VERSION 3.0.3
ENV HADOOP_URL http://mirror.bit.edu.cn/apache/hadoop/common/hadoop-${HADOOP_VERSION}/hadoop-${HADOOP_VERSION}.tar.gz

RUN set -x \
    && curl -fSL "$HADOOP_URL" -o /tmp/hadoop.tar.gz \
    && tar -xvf /tmp/hadoop.tar.gz -C /opt/ \
    && rm /tmp/hadoop.tar.gz*

RUN ln -s /opt/hadoop-${HADOOP_VERSION}/etc/hadoop /etc/hadoop
RUN mkdir /opt/hadoop-${HADOOP_VERSION}/logs

ENV HADOOP_PREFIX=/opt/hadoop-${HADOOP_VERSION}
ENV HADOOP_CONF_DIR=/etc/hadoop

ENV USER=root
ENV PATH $HADOOP_PREFIX/bin/:$PATH

ADD entry_point.sh /entry_point.sh
ADD hadoop_conf /hadoop_conf

RUN chmod a+x /entry_point.sh /hadoop_conf

ENTRYPOINT ["/entry_point.sh"]
