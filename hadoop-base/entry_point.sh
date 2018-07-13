#!/bin/bash

export CORE_CONF_fs_defaultFS=hdfs://hadoop.name:8020

function addProperty {
    local path=$1
    local name=$2
    local value=$3
    local property="<property><name>$name</name><value>${value}</value></property>"
    local escaped=$(echo $property | sed 's/\//\\\//g')
    sed -i "/<\/configuration>/ s/.*/${escaped}\n&/" $path 
}

function configEnv {
    local path=$1
    local module=$2
    local envPrefix=$3
    local keyValue
    local name
    local value

    echo "configuring $module"
    for c in `printenv | grep $envPrefix`; do
        keyValue=${c:${#envPrefix}}
        name=`echo $keyValue | cut -d'=' -f 1 | sed 's/___/-/g; s/__/_/g; s/_/./g'`
        value=`echo $keyValue | cut -d'=' -f 2
        echo "setting $name=$value"
        addProperty $path $name $value
    done
}

configEnv /etc/hadoop/core-site.xml core CORE_CONF_
configEnv /etc/hadoop/hdfs-site.xml hdfs HDFS_CONF_
configEnv /etc/hadoop/yarn-site.xml yarn YARN_CONF_
configEnv /etc/hadoop/mapred-site.xml mapred MAPRED_CONF_

exec $@