#!/bin/bash

filePrefix=file://
namedir=${HDFS_CONF_dfs_namenode_name_dir:${#filePrefix}}
echo $namedir

if [ ! -d $namedir ]; then
  echo "namenode name directory not found: $namedir"
  exit 2
fi

if [ "`ls -A $namedir`" == "" ]; then
  echo "formatting namenode name directory: $namedir"
  $HADOOP_HOME/bin/hdfs --config $HADOOP_CONF_DIR namenode -format 
fi

$HADOOP_HOME/bin/hdfs --config $HADOOP_CONF_DIR namenode