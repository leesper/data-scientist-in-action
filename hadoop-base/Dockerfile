FROM leesper/ubuntu-java

ENV HADOOP_VERSION 2.7.6
ENV HADOOP_URL http://mirror.bit.edu.cn/apache/hadoop/common/hadoop-${HADOOP_VERSION}/hadoop-${HADOOP_VERSION}.tar.gz

RUN set -x \
    && curl -fSL "$HADOOP_URL" -o /tmp/hadoop.tar.gz \
    && tar -xvf /tmp/hadoop.tar.gz -C /usr/local \
    && rm /tmp/hadoop.tar.gz*

ENV HADOOP_PREFIX=/usr/local/hadoop-${HADOOP_VERSION}
ENV HADOOP_HOME=${HADOOP_PREFIX}
ENV HADOOP_COMMON_HOME=${HADOOP_PREFIX}
ENV HADOOP_HDFS_HOME=${HADOOP_PREFIX}
ENV HADOOP_MAPRED_HOME=${HADOOP_PREFIX}
ENV HADOOP_YARN_HOME=${HADOOP_PREFIX}
ENV HADOOP_CONF_DIR=${HADOOP_PREFIX}/etc/hadoop
ENV YARN_CONF_DIR=${HADOOP_PREFIX}/etc/hadoop
ENV PATH=$HADOOP_HOME/bin/:$PATH
