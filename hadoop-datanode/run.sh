#!/bin/bash

filePrefix=file://
datadir=${HDFS_CONF_dfs_datanode_data_dir:${#filePrefix}}

if [ ! -d $datadir ]; then
  echo "datanode data directory not found: $datadir"
  exit 2
fi

$HADOOP_HOME/bin/hdfs --config $HADOOP_CONF_DIR datanode