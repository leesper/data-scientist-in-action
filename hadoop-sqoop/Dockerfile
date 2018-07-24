FROM leesper/hadoop-hbase

ENV SQOOP_URL=https://mirrors.tuna.tsinghua.edu.cn/apache/sqoop/1.4.7/sqoop-1.4.7.bin__hadoop-2.6.0.tar.gz

RUN set -x \
    && curl -fSL "$SQOOP_URL" -o /tmp/sqoop.tar.gz \
    && tar -xvf /tmp/sqoop.tar.gz -C /usr/local \
    && mv /usr/local/sqoop-1.4.7.bin__hadoop-2.6.0 /usr/local/sqoop \
    && rm -rf /tmp/sqoop.tar.gz* \
    && cp /usr/local/hive/lib/mysql-connector-java-5.1.46-bin.jar /usr/local/sqoop/lib

ENV SQOOP_HOME=/usr/local/sqoop
ENV PATH=$PATH:${SQOOP_HOME}/bin
ENV CLASSPATH=${CLASSPATH}:${SQOOP_HOME}/lib
