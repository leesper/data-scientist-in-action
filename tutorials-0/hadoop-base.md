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
5. dfs.replication：副本因子，每个数据块存多少个副本
6. dfs.webhdfs.enabled：是否启动NameNode和DataNode中的WebHDFS(REST API)
7. dfs.permissions.enabled：是否开启HDFS的文件权限检查
8. dfs.datanode.max.transfer.threads：DataNode传输数据的最大线程数

### 3. 配置文件mapred-site.xml

1. mapreduce.framework.name：执行MapReduce作业的运行时框架
2. mapreduce.jobhistory.address：MapReduce作业历史服务器主机名和端口
3. mapreduce.jobhistory.webapp.address：MapReduce作业历史服务器WebUI主机名和端口
4. mapreduce.jobhistory.done-dir：MR JobHistory Server管理的日志的存放位置
5. mapreduce.jobhistory.intermediate-done-dir：MapReduce作业产生的日志存放位置
6. mapreduce.map.log.level：Map任务日志级别
7. mapreduce.reduce.log.level：Reduce任务日志级别

### 4. 配置文件yarn-site.xml

1. yarn.nodemanager.aux-services：NodeManager上运行的附属服务，配置成mapreduce_shuffle才能运行MapReduce程序
2. yarn.nodemanager.remote-app-log-dir：应用程序日志聚集的HDFS目录
3. yarn.nodemanager.resource.memory-mb：表示该节点上YARN可使用的物理内存总量
4. yarn.nodemanager.resource.cpu-vcores：表示该节点上YARN可使用的虚拟CPU个数
5. yarn.nodemanager.pmem-check-enabled：是否启动一个线程检查每个任务正在使用的物理内存量
6. yarn.nodemanager.vmem-check-enabled：是否启动一个线程检查每个任务正在使用的虚拟内存量
7. yarn.resourcemanager.hostname：ResourceManager主机名
8. yarn.resourcemanager.address：ResourceManager对客户端暴露的地址，客户端通过该地址向RM提交应用程序，杀死应用程序等
9. yarn.resourcemanager.scheduler.address：ResourceManager对ApplicationMaster暴露的访问地址。ApplicationMaster通过该地址向RM申请资源、释放资源等
10. yarn.resourcemanager.resource-tracker.address：ResourceManager对NodeManager暴露的地址.。NodeManager通过该地址向RM汇报心跳，领取任务等
11. yarn.resourcemanager.admin.address：ResourceManager对管理员暴露的访问地址。管理员通过该地址向RM发送管理命令等
12. yarn.resourcemanager.webapp.address：ResourceManager对外Web UI地址，用户可通过该地址在浏览器中查看集群各类信息
13. yarn.scheduler.maximum-allocation-mb：单个容器可申请的最大内存资源量
14. yarn.scheduler.maximum-allocation-vcores：单个可申请的最大虚拟CPU个数
15. yarn.log-aggregation-enable：是否启用日志聚集功能查看container日志
16. yarn.log.server.url：Log server地址
17. yarn.timeline-service.enabled：是否开启timeline service服务
18. yarn.system-metrics-publisher.enabled：该设置控制RM是否发布YARN系统度量值到timeline server
19. yarn.timeline-service.generic-application-history.enabled：标示client是否通过timeline history-service查询通用的application数据
20. yarn.timeline-service.hostname：Timeline service web应用的主机名
21. yarn.app.mapreduce.am.staging-dir：作业启动后，Hadoop会将作业日志放在该目录下

## 编排并启动Hadoop集群

写好各种配置信息后，我们就可以正式开始编排并启动我们的Hadoop集群了。该集群由4个节点组成：一个NameNode节点，一个ResourceManager节点和两个DataNode节点，分别命名为：

1. master.namenode
2. master.resourcemanager
3. worker.datanode1
4. worker.datanode2

根据[Hadoop官方文档](http://hadoop.apache.org/docs/r2.8.4/hadoop-project-dist/hadoop-common/ClusterSetup.html)，一个典型的集群包含两个master：NameNode和ResourceManager，而DataNode一般和NodeManager部署在一台机器上。NameNode和DataNode作为HDFS的守护进程负责管理数据；而ResourceManager和NodeManager作为守护进程负责YARN的资源调度，将NodeManager和DataNode放在一起有助于实现“计算跟着数据走”，提高集群的运算效率。

docker编排文件如下：
```docker
version: "3"
services: 
  master.namenode:
    image: leesper/hadoop-base
    container_name: master.namenode
    volumes:
      - ./hadoop-base/conf:/usr/local/hadoop/etc/hadoop
    network_mode: zoo-net
    ports:
      - "50070:50070"
  master.resourcemanager:
    image: leesper/hadoop-base
    container_name: master.resourcemanager
    volumes: 
      - ./hadoop-base/conf:/usr/local/hadoop/etc/hadoop
    network_mode: zoo-net
    ports:
      - "8088:8088"
  worker.datanode1:
    image: leesper/hadoop-base
    container_name: worker.datanode1
    depends_on:
      - master.namenode
    volumes:
      - ./hadoop-base/conf:/usr/local/hadoop/etc/hadoop
    network_mode: zoo-net
  worker.datanode2:
    image: leesper/hadoop-base
    container_name: worker.datanode2
    depends_on:
      - master.namenode
    volumes:
      - ./hadoop-base/conf:/usr/local/hadoop/etc/hadoop
    network_mode: zoo-net

```

这里有几个值得注意的地方。首先，通过设置container_name指定容器名，Docker能够自动建立IP到container_name的解析，集群中的各个节点就能通过它互相访问，所以编排文件中写的与配置文件是保持一致的；其次，volumes指定了将宿主机中的某个文件夹挂载到虚拟机的某个文件夹，这相当于建立了映射。所以这里只要将当前目录下的配置文件映射到虚拟机中对应的文件夹，就相当于建立了配置文件，就不用在构建镜像时进行配置了，当然，你也可以通过写shell脚本的方式进行配置；第三，network_mode指定了我们之前创建的Docker网络：zoo-net，这样一来所有的容器就处于同一个网段中；最后，depends_on指定了容器之间的依赖关系，编排文件启动Docker集群时会将被依赖的容器优先启动起来。

## 参考文献

[大数据技术原理与应用：概念、存储、处理、分析与应用（第2版）](https://www.amazon.cn/dp/B06X1DYSBS/ref=sr_1_1?ie=UTF8&qid=1534749671&sr=8-1&keywords=%E5%A4%A7%E6%95%B0%E6%8D%AE%E6%8A%80%E6%9C%AF%E5%8E%9F%E7%90%86%E4%B8%8E%E5%BA%94%E7%94%A8%EF%BC%9A%E6%A6%82%E5%BF%B5%E3%80%81%E5%AD%98%E5%82%A8%E3%80%81%E5%A4%84%E7%90%86%E3%80%81%E5%88%86%E6%9E%90%E4%B8%8E%E5%BA%94%E7%94%A8%EF%BC%88%E7%AC%AC2%E7%89%88%EF%BC%89)