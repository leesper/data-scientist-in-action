FROM leesper/hadoop-base:3.0.3

ENV YARN_CONF_yarn_app_mapreduce_am_staging___dir=/mr

RUN mkdir /mr
VOLUME /mr

ADD run.sh /run.sh
RUN chmod a+x /run.sh

CMD ["/run.sh"]
