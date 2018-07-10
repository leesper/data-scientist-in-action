FROM leesper/hadoop-base:3.0.3

ENV HDFS_CONF_dfs_namenode_name_dir=file:///data/hadoop/name
RUN mkdir -p /data/hadoop/name
VOLUME /data/hadoop/name

ADD run.sh /run.sh
RUN chmod a+x /run.sh

CMD ["/run.sh"]