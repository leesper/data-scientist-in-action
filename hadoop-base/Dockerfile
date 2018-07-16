FROM ubuntu:18.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends openjdk-8-jdk openssh-server net-tools curl \
    && rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64

ENV HADOOP_VERSION 2.7.3
ENV HADOOP_URL http://mirror.bit.edu.cn/apache/hadoop/common/hadoop-${HADOOP_VERSION}/hadoop-${HADOOP_VERSION}.tar.gz

RUN set -x \
    && curl -fSL "$HADOOP_URL" -o /tmp/hadoop.tar.gz \
    && tar -xvf /tmp/hadoop.tar.gz -C /usr/local \
    && rm /tmp/hadoop.tar.gz* \
    && ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa \
    && cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys \
    && chmod 0600 ~/.ssh/authorized_keys

ENV HADOOP_PREFIX=/usr/local/hadoop-${HADOOP_VERSION}
ENV HADOOP_HOME=${HADOOP_PREFIX}
ENV HADOOP_COMMON_HOME=${HADOOP_PREFIX}
ENV HADOOP_HDFS_HOME=${HADOOP_PREFIX}
ENV HADOOP_MAPRED_HOME=${HADOOP_PREFIX}
ENV HADOOP_YARN_HOME=${HADOOP_PREFIX}
ENV HADOOP_CONF_DIR=${HADOOP_PREFIX}/etc/hadoop
ENV YARN_CONF_DIR=${HADOOP_PREFIX}/etc/hadoop
ENV PATH=$HADOOP_HOME/bin/:$PATH
ENV USER=root
