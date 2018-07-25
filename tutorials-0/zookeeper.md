# 实战任务0：利用Docker搭建ZooKeeper集群

大家好！“搭建完全分布式的大数据集群”第一关要完成的实战任务是利用Docker搭建ZooKeeper集群。作为实验指导，本文首先会对ZooKeeper及其应用场景进行一个简单的介绍；然后手把手教大家构建基础操作系统镜像，这将会是后面很多工作的起点，在这个过程中也会给大家介绍一点编写Dockerfile的技巧；紧接着，我们将趁热打铁，使用刚学会的技巧来构建ZK镜像，在构建镜像的过程中会介绍与ZK有关的配置选项；最后，通过学习编排文件的编写技巧，我们将掌握在本地启动ZK集群的技能。

## ZooKeeper介绍

ZK是一套中心化的服务，它提供高可用的，与分布式协同工作相关的服务，例如统一命名服务，配置管理，集群管理和分布式锁。它具备高吞吐量，低延迟和高可用等特点，通过Leader Election机制，ZK集群避免了传统意义上因为单点故障引起整个服务不可用的问题。

## ZooKeeper典型应用场景

ZK通过观察者模式来实现不同服务之间的分布式协作。它管理和维护大家都关心的数据。当数据发生变化时，注册在它上面的观察者就会收到通知并做出相应的反应。下面是一些典型的应用场景。

### 1. 统一命名服务

分布式应用中，通常需要有一套完整的命名规则，既能够产生唯一的名称又便于人识别和记住，通常情况下用树形的名称结构是一个理想的选择，树形的名称结构是一个不会重复的有层次的目录结构，就像数据库中产生一个唯一的数字主键一样。Name Service已经是ZK内置的功能，只要调用API就能实现。如调用create接口就可以很容易创建一个目录节点。

### 2. 配置管理

配置的管理在分布式服务中很常见，例如同一个服务运行在多台机器上，但是某些配置项是相同的，如果要修改这些相同的配置项，那么就必须同时修改每台机器，这样非常麻烦而且容易出错。这样的配置信息完全可以交给ZK来管理，将配置信息保存在ZK的某个目录节点中，然后将所有需要修改的应用机器监控配置信息的状态，一旦配置信息发生变化，每台应用机器就会收到ZK的通知，然后从ZK获取新的配置信息应用到系统中。

### 3. 集群管理

ZK能够很容易的实现集群管理的功能。例如多台机器组成的服务器集群，那么必须要有一个“总管”知道当前集群中每台机器的服务状态，一旦添加了新的机器进行扩容或者有机器出现故障不能提供服务，集群中其它机器必须知道，从而做出调整重新分配服务策略。ZK不仅能够帮你维护当前的集群中机器的服务状态，而且能够帮你选出一个“总管”，让这个总管来管理集群，这就是ZK的另一个功能Leader Election。它们的实现方式都是在ZK上创建一个EPHEMERAL类型的目录节点，然后每个服务在它们创建目录节点的父目录节点上调用getChildren(String path, boolean watch) 方法并设置watch为true，由于是EPHEMERAL目录节点，当创建它的机器出现故障挂了，这个目录节点也随之被删除，所以Children将会变化，这时getChildren上的Watch将会被调用，所以其它机器就知道已经有某台机器故障了。新增机器进行扩容时也是同样的原理。

ZK选出一个Master Server的Leader Election机制也是类似的原理：每台机器会创建一个EPHEMERAL目录节点，不同的是它还是一个SEQUENTIAL目录节点，所以它是个EPHEMERAL_SEQUENTIAL目录节点。之所以它是EPHEMERAL_SEQUENTIAL目录节点，是因为我们可以给每台机器编号，我们可以选择当前是最小编号的机器为Master，假如这个最小编号的机器故障了，由于是EPHEMERAL节点，挂掉的机器对应的节点也被删除，所以当前的节点列表中又出现一个最小编号的节点，我们就选择这个节点为当前Master。这样就实现了动态选择Master，避免了单点故障的问题。

### 4. 分布式锁

分布式锁在同一个进程中很容易实现，但是在跨进程或者在不同服务之间就不好实现了。ZK很容易实现这个功能，实现方式是需要获得锁的服务创建一个EPHEMERAL_SEQUENTIAL目录节点，然后调用getChildren方法获取当前的目录节点列表中最小的目录节点，然后判断是不是就是自己创建的目录节点，如果正是自己创建的，那么它就获得了这个锁，如果不是那么它就调用exists(String path, boolean watch)方法并监控ZK上目录节点列表的变化，一直到自己创建的节点是列表中最小编号的目录节点，从而获得锁，释放锁很简单，只要删除前面它自己所创建的目录节点就行了。

## 构建操作系统基础镜像

首先，我们需要构建一个操作系统基础镜像，选用任何Linux发行版都可以，我本人比较熟悉Ubuntu，所以这里就使用官方的Ubuntu18.04镜像。后面几乎所有的镜像都是基于该镜像构建的，所以在这个镜像中还会配置和安装好OpenSSH和Java8，设置好相关环境变量，配置好SSH免密登录，最后前台启动sshd服务，这样就能保证该容器长时间运行，而不是一启动就退出了。ubuntu-java目录下完整的Dockerfile如下：

```docker
FROM ubuntu:18.04

RUN apt-get update \ # 更新apt仓库，安装软件包
    && apt-get install -y --no-install-recommends openjdk-8-jdk openssh-server net-tools curl \
    && rm -rf /var/lib/apt/lists/* \ # 删除下载的安装信息，给镜像瘦身
    && set -x \
    && ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa \ # 配置SSH互信免密
    && cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys \
    && chmod 0600 ~/.ssh/authorized_keys \
    && mkdir -p /var/run/sshd

ENV USER=root
ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64 # Java环境变量
ENV PATH=${JAVA_HOME}/bin:$PATH # 更新系统PATH环境变量

EXPOSE 22 # 运行时容器将监听该端口（TCP）
CMD ["/usr/sbin/sshd", "-D"]
```

编写Dockerfile时有三个技巧为构建出来的镜像瘦身。首先，选用较小的操作系统镜像作为基础，这里我选择的官方Ubuntu镜像只有81.2MB，当然这里还可以选取更小的操作系统镜像；其次，只安装必要的软件，安装完毕后将文件删除和清理掉；最后，因为Docker镜像的构建都是一层一层的，层数越多体积就越大，所以还有一个技巧就是尽量将RUN指令写在一起，利用“\”和“&&”符号把多个命令写到一个RUN指令中，减少镜像构建所需要的层数。上面编写的Dockerfile很好地体现了这样的原则。

## 构建ZooKeeper镜像

首先，我们要准备一个单独的目录zookeeper，在该目录下有一个conf子目录存放了所有的配置文件。这里准备配置一个由3个节点组成的完全分布式的ZK集群，所以该目录下还有三个子目录zk1，zk2和zk3，每个目录下面有个myid文件，分别包含3个节点的ID（1，2和3），配置文件如下：

```
tickTime=2000  # 单位毫秒，通常用作心跳间隔时间，最小会话超时时间为该时间两倍
dataDir=/opt/data  # 存储in-memory数据库快照的路径
dataLogDir=/opt/log  # 存储事务日志的路径
clientPort=2181  # 监听客户端连接的端口
initLimit=10  # timeouts ZK uses to limit the length of time the ZK servers in quorum have to connect to a leader
syncLimit=5  # limits how far out of date a server can be from a leader

# peers use former port to connect to other peers
# the latter port is used for leader election
server.1=zk1:2888:3888
server.2=zk2:2888:3888
server.3=zk3:2888:3888
```

下面在ubuntu-java镜像的基础上，构建zookeeper镜像，启动后的ZK集群有3个节点，每个节点都会前台运行zkServer.sh，Dockerfile如下：

```docker
FROM leesper/ubuntu-java

ENV ZOOKEEPER_VERSION=3.4.10
ENV ZOOKEEPER_URL=http://mirrors.shu.edu.cn/apache/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/zookeeper-${ZOOKEEPER_VERSION}.tar.gz

RUN set -x \
    && curl -fSL "${ZOOKEEPER_URL}" -o /tmp/zookeeper.tar.gz \
    && tar -xvf /tmp/zookeeper.tar.gz -C /usr/local \
    && mv /usr/local/zookeeper-${ZOOKEEPER_VERSION} /usr/local/zookeeper \
    && rm -rf /tmp/zookeeper.tar.gz* \
    && mkdir -p /opt/data \
    && mkdir -p /opt/log

ENV ZOO_HOME=/usr/local/zookeeper
ENV PATH=$PATH:${ZOO_HOME}/bin
ENV TZ="Asia/Shanghai"
EXPOSE 2181 2888 3888
CMD ["zkServer.sh", "start-foreground"]
```

## 编排并启动ZooKeeper集群

基础镜像构建完毕后，可以开始编排ZK集群的启动了，对应的docker-compose-zk.yml编排文件如下所示：

```docker
version: "3"
services: 
  zk1:
    image: leesper/zookeeper
    container_name: zk1  # 指定容器名后docker会自动将这个名字与其IP地址建立域名解析关系
    network_mode: zoo-net  # ZK集群运行于我们事先创建好的zoo-net容器网络中
    volumes:  
      # 将宿主机的zookeeper/zk1挂载到容器的/opt/data
      - ./zookeeper/zk1:/opt/data  
      # 将宿主机的zookeeper/conf挂载到容器的/usr/local/zookeeper/conf
      - ./zookeeper/conf:/usr/local/zookeeper/conf  
    ports:  # 将容器的2181端口映射到宿主机2181端口
      - "2181:2181"
    expose:  # 容器要监听的端口
      - "2888"
      - "3888"
  zk2:
    image: leesper/zookeeper
    container_name: zk2
    network_mode: zoo-net
    volumes:
      - ./zookeeper/zk2:/opt/data 
      - ./zookeeper/conf:/usr/local/zookeeper/conf
    ports:
      - "2182:2181"
    expose:
      - "2888"
      - "3888"
  zk3:
    image: leesper/zookeeper
    container_name: zk3
    network_mode: zoo-net
    volumes:
      - ./zookeeper/zk3:/opt/data
      - ./zookeeper/conf:/usr/local/zookeeper/conf
    ports:
      - "2183:2181"
    expose:
      - "2888"
      - "3888"
```

通过`docker-compose -f docker-compose-zk.yml up -d`启动ZK集群，可根据需要在编排文件中增减集群数量，但注意同时要增减myid配置，并修改相关配置文件。通过`docker ps`可看到当前已启动的容器。

## 参考文献

1. [分布式服务框架 Zookeeper -- 管理分布式环境中的数据](https://www.ibm.com/developerworks/cn/opensource/os-cn-zookeeper/index.html)

2. [Apache ZooKeeper](https://zookeeper.apache.org/)

3. [ZooKeeper Wiki](https://cwiki.apache.org/confluence/display/ZOOKEEPER/Index)