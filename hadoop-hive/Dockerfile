FROM leesper/hadoop-base

ENV HIVE_VERSION=1.2.2
ENV HIVE_URL=https://mirrors.tuna.tsinghua.edu.cn/apache/hive/hive-${HIVE_VERSION}/apache-hive-${HIVE_VERSION}-bin.tar.gz
ENV MYSQL_CONNECTOR_URL=https://cdn.mysql.com//Downloads/Connector-J/mysql-connector-java-5.1.46.tar.gz

RUN set -x \
    && curl -fSL "$HIVE_URL" -o /tmp/hive.tar.gz \
    && tar -xvf /tmp/hive.tar.gz -C /usr/local \
    && mv /usr/local/apache-hive-${HIVE_VERSION}-bin /usr/local/hive \
    && rm -rf /tmp/hive.tar.gz* \
    && curl -fSL "$MYSQL_CONNECTOR_URL" -o /tmp/mysql-connector-java.tar.gz \
    && tar -xvf /tmp/mysql-connector-java.tar.gz -C /usr/local \
    && mv /usr/local/mysql-connector-java-5.1.46/mysql-connector-java-5.1.46-bin.jar ${HIVE_HOME}/lib \
    && rm -rf /tmp/mysql-connector-java.tar.gz* \
    && mkdir -p /usr/hive/warehouse \
    && mkdir -p /usr/hive/log

ENV HIVE_HOME=/usr/local/hive
ENV PATH=$PATH:${HIVE_HOME}/bin:.
