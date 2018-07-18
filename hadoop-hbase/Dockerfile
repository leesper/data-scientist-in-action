FROM leesper/hadoop-spark

ENV HBASE_VERSION=2.0.1
ENV HBASE_URL=http://mirrors.shu.edu.cn/apache/hbase/${HBASE_VERSION}/hbase-${HBASE_VERSION}-bin.tar.gz

RUN set -x \
    && curl -fSL "$HBASE_URL" -o /tmp/hbase.tar.gz \
    && tar -xvf /tmp/hbase.tar.gz -C /usr/local \
    && mv /usr/local/hbase-${HBASE_VERSION} /usr/local/hbase \
    && rm -rf /tmp/hbase.tar.gz*

ENV HBASE_HOME=/usr/local/hbase
ENV PATH=$PATH:${HBASE_HOME}/bin

# HDFS ports
EXPOSE 9000 50010 50020 50070 50075 50090

# See https://issues.apache.org/jira/browse/HDFS-9427
# Updates the HDFS default HTTP/RPC ports to non-ephemeral ports. Changes are:
# NameNode ports: 50470 --> 9871, 50070 --> 9870, 8020 --> 9820
# Secondary NN ports: 50091 --> 9869, 50090 --> 9868
# Datanode ports: 50020 --> 9867, 50010 --> 9866, 50475 --> 9865, 50075 --> 9864
EXPOSE 9871 9870 9820 9869 9868 9867 9866 9865 9864

# MAPRED ports
EXPOSE 19888

# YARN ports
EXPOSE 8030 8031 8032 8033 8040 8042 8088 8188

# Other ports
EXPOSE 49707 2122