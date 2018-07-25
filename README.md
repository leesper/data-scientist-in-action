# data-scientist-in-action 行动中的数据科学家

[![GitHub forks](https://img.shields.io/github/forks/leesper/data-scientist-in-action.svg)](https://github.com/leesper/data-scientist-in-action/network) [![GitHub stars](https://img.shields.io/github/stars/leesper/data-scientist-in-action.svg)](https://github.com/leesper/data-scientist-in-action/stargazers) [![GitHub license](https://img.shields.io/github/license/leesper/data-scientist-in-action.svg)](https://github.com/leesper/data-scientist-in-action/blob/master/LICENSE)

comprehensive projects for data engineering and analysis

大数据工程和分析综合案例

## Features

通过本项目，你可以学习到：

1. 如何编写Dockerfile来构建各种常用Hadoop组件容器
2. 如何基于Docker的编排技术来构建完全分布式的大数据容器集群
3. 通过自己搭建的Hadoop容器集群对[Kaggle泰坦尼克数据集](https://www.kaggle.com/c/titanic)进行数据分析

## Environments

* Ubuntu 18.04
* OpenJDK Java 8
* Zookeeper 3.4.10
* Hadoop 2.7.6
* mysql 5.6.40
* Hive 1.2.2
* Spark 2.3.1
* Hbase 2.0.1
* Sqoop 1.4.7

## Images

1. leesper/ubuntu-java：Java8 + OpenSSH，基础操作系统镜像
2. leesper/zookeeper：基于leesper/ubuntu-java构建，用于启动ZooKeeper集群
3. leesper/mysql：基于官方镜像构建，用于启动MySQL容器提供给集群使用
4. leesper/hadoop-base：基于leesper/ubuntu-java构建，用于启动基础Hadoop集群
6. leesper/hadoop-hive：基于leesper/hadoop-base构建，包含Hadoop和Hive，用于启动带Hive的Hadoop集群
7. leesper/hadoop-spark：基于leesper/hadoop-hive构建，包含Hadoop，Hive和Spark，用于启动Hadoop+Spark集群
8. leesper/hadoop-hbase：基于leesper/hadoop-spark构建，包含Hadoop，Hive，Spark和HBase，用于启动Hadoop+Spark+HBase集群
9. leesper/hadoop-sqoop：基于leesper/hadoop-hbase构建，包含Hadoop，Hive，Spark，HBase和Sqoop，提供Sqoop工具

## Requirements

本项目在以下环境中测试通过：

* Ubuntu Linux 18.04
* Docker 18.03.1-ce
* Docker Compose 1.17.1

## Tutorials 0：搭建完全分布式的大数据集群

### 1. 构建/拉取镜像

所有的镜像都已构建完毕并上传到DockerHub，可以通过[这里](https://hub.docker.com/u/leesper/)拉取所需镜像，也可以通过`sh docker_pull.sh`一次性拉取所有构建好的镜像。

若想自己构建镜像，请进入相应目录并运行对应的`docker build`命令，以ubuntu-java为例：

```
cd ubuntu-java/

# 注意末尾小数点不可省略（代表当前目录）
docker build -t=xxx/ubuntu-java .
```

### 2. 创建大数据集群网络

通过命令：`docker network create zoo-net`创建名称为zoo-net的容器网络，**注意**：若指定了不同的名称，需要修改对应的Docker编排文件。

### 3. 搭建ZooKeeper集群

[实战项目0：利用Docker搭建ZooKeeper集群](./tutorials-0/zookeeper.md)

### 4. 启动MySQL容器

[实战项目1：利用Docker搭建MySQL Server容器](./tutorials-0/mysql.md)

### 5. 启动基础Hadoop集群

**TODO** 增加对组件的描述

```
$ docker-compose -f docker-compose-hadoop.yml up -d
```

启动集群，格式化NameNode

```
docker exec -it master.namenode /bin/bash
hdfs namenode -format
```

在NameNode上启动HDFS
```
cd /usr/local/hadoop/sbin
./start-dfs.sh
```
然后在ResourceManager上启动YARN和MapReduce JobHistory Server

```
cd /usr/local/hadoop/sbin
./start-yarn.sh
./mr-jobhistory-daemon.sh --config $HADOOP_CONF_DIR start historyserver
```

访问http://localhost:50070，看集群是否启动成功，也可以通过jps命令查看进程。

### 6. 启动Hadoop+Hive集群

需要依赖mysql容器

```
$ docker-compose -f docker-compose-hive.yml up -d
```

需要先启动Hadoop基础集群，操作同上

可以在任意一个节点上启动Hive Server，启动之后可以以编程的方式访问Hive中的数据

```
hive --service hiveserver2 # 这里是在NameNode节点上运行的，默认前台运行，若需要后台运行，请加&
```

以交互式方式访问Hive

```
hive
```

### 7. 启动Hadoop+Hive+Spark集群

需要依赖mysql容器

```
$ docker-compose -f docker-compose-spark.yml up -d
```

需要先启动Hadoop基础集群，操作同上

NameNode上，启动Spark集群

```
/usr/local/spark/sbin/start-all.sh
```

Standalone模式下，使用Spark自带样例中的计算Pi的应用来验证一下

```
spark-submit --master spark://master.namenode:7077 --class org.apache.spark.examples.SparkPi /usr/local/spark/examples/jars/spark-examples_2.11-2.3.1.jar 1000
```

计算结果输出如下

```
18/07/24 03:57:16 WARN NativeCodeLoader: Unable to load native-hadoop library for your platform... using builtin-java classes where applicable
Pi is roughly 3.1413606314136064
```

YARN模式下，使用Spark自带样例中的计算Pi的应用来验证一下
```
spark-submit --master yarn --class org.apache.spark.examples.SparkPi /usr/local/spark/examples/jars/spark-examples_2.11-2.3.1.jar 1000
```

计算结果输出如下
```
18/07/24 04:00:36 WARN NativeCodeLoader: Unable to load native-hadoop library for your platform... using builtin-java classes where applicable
18/07/24 04:00:41 WARN Client: Neither spark.yarn.jars nor spark.yarn.archive is set, falling back to uploading libraries under SPARK_HOME.
Pi is roughly 3.1417248314172483
```

### 8. 启动Hadoop+Hive+Spark+Hbase集群

```
docker-compose -f docker-compose-hbase.yml up -d
```

需要先启动Hadoop基础集群和Spark集群，操作同上

在NameNode上启动HBase集群

```
/usr/local/hbase/bin/start-hbase.sh
```

### 9. 启动Hadoop+Hive+Spark+Hbase+Sqoop集群

```
docker-compose -f docker-compose-sqoop.yml up -d
```

需要先启动Hadoop基础集群，Spark集群和HBase集群，操作同上

在任意一个节点上（比如DataNode），测试Sqoop与MySQL之间的连接是否成功：

```
sqoop list-databases --connect jdbc:mysql://mysql:3306/ --username root -Proot
```

## Support

Tell people where to get help, such as issue tracker, chat room, an email, etc.

## Development
For people who want to make changes to your project, it's helpful to have some documentation on how to get started. Such as how to setup an environment. For example:

```
$ virtualenv foobar
$ . foobar/bin/activate
$ pip install -e .
```

## Authors and acknowledgment

Show your appreciation to those who have contributed to the project.

## Changelog

A record of all notable changes made to a project.


注意docker-compose-hadoop.yml、docker-compose-hive.yml、docker-compose-spark.yml和docker-compose-hbase.yml不要一起启动，后面模板中是包含了前一个的所有配置
