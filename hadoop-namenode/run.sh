#!/bin/bash

namedir=`echo $HDFS_CONF_dfs.namenode.name.dir | perl -pe 's#file://##'`
if [ ! -d $namedir ]; then
  echo "namenode name directory not found: $namedir"
  exit 2
fi

if [ "`ls -A $namedir`" == "" ]; then
  echo "formatting namenode name directory: $namedir"
  $HADOOP_PREFIX/bin/hdfs --config $HADOOP_CONF_DIR namenode -format 
fi

$HADOOP_PREFIX/bin/hdfs --config $HADOOP_CONF_DIR namenode