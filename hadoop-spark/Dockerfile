FROM leesper/hadoop-hive

ENV SPARK_VERSION=2.3.1
ENV SPARK_URL=https://archive.apache.org/dist/spark/spark-${SPARK_VERSION}/spark-${SPARK_VERSION}-bin-without-hadoop.tgz

RUN set -x \
    && curl -fSL "$SPARK_URL" -o /tmp/spark.tgz \
    && tar -xvf /tmp/spark.tgz -C /usr/local \
    && mv /usr/local/spark-${SPARK_VERSION}-bin-without-hadoop /usr/local/spark \
    && rm -rf /tmp/spark.tgz*

ENV SPARK_HOME=/usr/local/spark
ENV PATH=$PATH:${SPARK_HOME}/bin:${SPARK_HOME}/sbin:.