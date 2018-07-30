# 实战任务2：利用Docker搭建ZooKeeper集群

完成了ZK和MySQL等容器的构建和编排之后，下面我们可以正式开始构建基础的Hadoop集群了。本次实战任务只编排和启动基础的Hadoop集群，也就是说，只包含HDFS，YARN和MapReduce这些核心组件。本关任务也是后面实战任务的基础，因为后面我们将以这里构建好的镜像为基础，逐步的添加一些其他的大数据组件：

* 实战任务3将在Hadoop基础集群的基础上增加Hive组件；
* 实战任务4将构建和编排Hadoop+Hive+Spark集群；
* 实战任务5将构建和编排Hadoop+Hive+Spark+HBase集群；
* 实战任务6将构建和编排Hadoop+Hive+Spark+HBase+Sqoop集群；

下面我们来简单介绍一下Hadoop核心组件。

## Apache Hadoop核心组件

对分布式系统的研究主要集中在两个方面：计算和存储。升级改良之后的Hadoop2主要包含了3个组件：HDFS，YARN和MapReduce，HDFS用于存储，而YARN和MapReduce则与计算相关。

### 1. HDFS

Hadoop分布式文件系统HDFS是谷歌分布式文件系统GFS的开源实现，通过网络实现文件在多台机器上的分布式存储，满足大规模数据存储的需求，这种规模和量级是单机绝对无法满足的。HDFS是按照主从结构设计的。主节点被称为名称节点NameNode，负责管理分布式文件系统的命名空间。两个核心数据结构为FsImage和EditLog。FsImage维护文件系统树，文件树中的所有文件和文件夹的元数据；EditLog则按顺序记录了所有对文件系统的操作。NameNode在启动时都会将FsImage载入内存，然后执行EditLog中的操作使得FsImage保持最新。数据节点DataNode是分布式文件系统中的工作节点，负责数据的存取，会根据NameNode节点的调度来进行数据的存储和检索，并定期向NameNode发送自己所存储的块列表，每个DataNode中的数据会被保存在各自节点的本地Linux文件系统中。

HDFS中还有一个SecondaryNameNode的设计。首先它每隔一段时间会和NameNode通信，完成EditLog和FsIma的合并操作，减少EditLog的大小，缩短NameNode节点重启时间；其次，作为NameNode检查点，因为其合并得到的FsImage可作为NameNode节点故障时恢复元数据的重要途径，虽然有一部分元数据仍然会丢失。

### 2. YARN

YARN是新一代的资源调度管理框架，同样采取的是主从结构设计。其中ResourceManager负责资源管理；ApplicationMaster负责任务调度和监控；NodeManager负责执行具体任务。它是一个纯粹的资源管理调度框架，甚至Spark分布式计算框架都可配置为运行在YARN之上。ResourceManager是YARN的主节点，包含Scheduler和ApplicationManager两个核心组件。Scheduler接收来自ApplicationMaster的应用程序资源请求，将集群中的资源以“容器”（动态资源分配单位，包含一定数量的CPU，内存和磁盘等资源）形式分配；ApplicationManager负责系统中所有应用程序的管理工作，包括应用程序提交，与Scheduler协商资源以启动ApplicationMaster，监控ApplicationMaster运行状态等。

用户提交的应用程序作业会被拆分成多个任务分布式执行。ResourceManager接收用户提交的作业，为其启动一个ApplicationMaster。ApplicationMaster会与ResourceManager协商获取以容器形式分配的资源，活得的资源会进一步分配给内部多个任务，实现资源二次分配。ApplicationMana会与NodeManager保持交互通信进行应用程序的启动，运行，监控和停止，在任务发生失败时进行恢复（重新申请资源并重启任务），并定时发送心跳报告资源使用情况和应用进度信息，作业完成时Application向ResourceManager注销容器，完成执行。

NodeManager驻留在YARN集群上每个节点的代理，负责每个节点的生命周期管理和监控资源使用情况。它会用心跳定时向ResourceManager报告各种情况，也会接收来自Applicat的启动/停止请求。实际部署的时候，ApplicationMaster和NodeManager一般会跟DataNode部署在一起，ResourceManager一般会跟NameNode部署在一起。

### 3. MapReduce

MapReduce是一种编程模型，它将运行于大规模分布式集群中的并行计算高度抽象为Map和Reduce两个过程，完成海量数据的处理，由于我们之后会采用更加强大的Spark框架，所以这里就不多做介绍了。

## 构建Hadoop镜像

构建Hadoop镜像比较简单，只要将对应的Hadoop压缩包下载下来，解压缩放到相应位置并设置好常用的环境变量即可，以下是Dockerfile文件：

```
FROM leesper/ubuntu-java

ENV HADOOP_VERSION 2.7.6
ENV HADOOP_URL http://mirror.bit.edu.cn/apache/hadoop/common/hadoop-${HADOOP_VERSION}/hadoop-${HADOOP_VERSION}.tar.gz

RUN set -x \
    && curl -fSL "$HADOOP_URL" -o /tmp/hadoop.tar.gz \
    && tar -xvf /tmp/hadoop.tar.gz -C /usr/local \
    && mv /usr/local/hadoop-${HADOOP_VERSION} /usr/local/hadoop \
    && rm -rf /tmp/hadoop.tar.gz*

ENV HADOOP_PREFIX=/usr/local/hadoop
ENV HADOOP_HOME=${HADOOP_PREFIX}
ENV HADOOP_COMMON_HOME=${HADOOP_PREFIX}
ENV HADOOP_HDFS_HOME=${HADOOP_PREFIX}
ENV HADOOP_MAPRED_HOME=${HADOOP_PREFIX}
ENV HADOOP_YARN_HOME=${HADOOP_PREFIX}
ENV HADOOP_CONF_DIR=${HADOOP_PREFIX}/etc/hadoop
ENV YARN_CONF_DIR=${HADOOP_PREFIX}/etc/hadoop
ENV PATH=$HADOOP_HOME/bin/:$PATH
```
关键在于core-site，hdfs-site，mapred-site和yarn-site四大配置文件，下面简要介绍一些常用的。

### 1. 配置文件core-site.xml

1. fs.defaultFS：默认文件系统的URL
2. hadoop.tmp.dir：存放临时文件的目录
3. fs.trash.interval：多少分钟后删除checkpoint数据

更多配置项及其默认值请参考[core-default.xml](http://hadoop.apache.org/docs/r2.7.6/hadoop-project-dist/hadoop-common/core-default.xml)文件。

### 2. 配置文件hdfs-site.xml

1. dfs.namenode.name.dir：NameNode在本地文件系统中存放fsimage的路径
2. dfs.datanode.data.dir：DataNode在本地文件系统中存放数据块block的路径
3. dfs.namenode.checkpoint.dir：SecondaryNameNode在本地文件系统中存放临时镜像的路径
4. dfs.namenode.secondary.http-address：SecondaryNameNode的HTTP服务器地址和端口
5. dfs.replication：
6. dfs.webhdfs.enabled
7. dfs.permissions
8. dfs.datanode.max.transfer.threads

### 3. 配置文件mapred-site.xml
### 4. 配置文件yarn-site.xml

## 编排并启动Hadoop集群

## 参考文献