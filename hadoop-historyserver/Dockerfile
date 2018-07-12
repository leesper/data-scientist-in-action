FROM leesper/hadoop-base:3.0.3

ENV YARN_CONF_yarn_timeline___service_leveldb___timeline___store_path=/data/hadoop/tmp
RUN mkdir -p /data/hadoop/tmp
VOLUME /data/hadoop/tmp

ADD run.sh /run.sh
RUN chmod a+x /run.sh

CMD ["/run.sh"]