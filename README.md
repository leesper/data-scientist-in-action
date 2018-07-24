# data-scientist-in-action 行动中的数据科学家
comprehensive projects for data engineering and analysis
大数据工程和分析综合案例

## 系统环境

* 操作系统：Ubuntu 18.04
* Java 8
* Zookeeper：3.4.10
* Hadoop：2.7.6
* mysql：5.6
* Hive：1.2.2
* Spark：2.3.1
* Hbase：2.0.1

## 镜像介绍

* leesper/ubuntu-java：openssh+Java 8 基础镜像 
* leesper/zookeeper：基于leesper/ubuntu-java构建，用于启动ZooKeeper集群
* leesper/hadoop-base：基于leesper/ubuntu-java构建，用于启动基础Hadoop集群
* leesper/hadoop-hive：基于leesper/hadoop-base镜像构建，用于启动带Hive组件的Hadoop集群
* leesper/hadoop-spark：基于leesper/hadoop-hive镜像构建，用于启动带Hive和Spark组件的Hadoop集群
* leesper/hadoop-hbase：基于leesper/hadoop-spark镜像构建，用于启动带Hive，Spark和HBase的Hadoop集群

## Quick Start

### 1. 构建/拉取镜像

可通过`sh docker_build.sh`命令一次性构建所有要用到的镜像，也可以通过`sh docker_pull`从[Docker Hub](https://hub.docker.com/u/leesper/)拉取已构建好的镜像。

### 2. 创建大数据集群网络

可通过`docker network create zoo-net`创建名称为zoo-net的容器网络，若指定了不同的名称，请注意修改对应的Docker编排文件。

### 3. 启动ZooKeeper集群

```
$ docker-compose -f docker-compose-zk.yml up -d
```

根据需要可在docker compose文件中增减集群数量，注意同时要增减myid配置

### 4. 启动MySQL容器

若仅仅想使用步骤5中的基础Hadoop集群，可省略此步。

```
$ docker-compose -f docker-compose-mysql.yml up -d
```

修改密码和配置远程访问mysql：

```
docker exec -it mysql /bin/bash
mysql -u root -proot
# 进入名为mysql的数据库
mysql> use mysql;
# 授权远程访问
mysql> GRANT ALL PRIVILEGES ON *.* TO root@"%" IDENTIFIED BY "root" WITH GRANT OPTION;
mysql> FLUSH PRIVILEGES;

# 配置字符集，解决后面Hive建表报错 #TODO
# FAILED: Execution Error, return code 1 from org.apache.hadoop.hive.ql.exec.DDLTask. MetaException(message:For direct MetaStore DB connections, we don't support retries at the client level.)

mysql> ALTER DATABASE hive character set latin1;
```

MySQL容器配置完成

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


注意docker-compose-hadoop.yml、docker-compose-hive.yml、docker-compose-spark.yml和docker-compose-hbase.yml不要一起启动，后面模板中是包含了前一个的所有配置
