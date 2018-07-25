# 实战项目1：利用Docker搭建MySQL Server容器

完成了对ZooKeeper的搭建后，下面我们来构建一个独立运行MySQL Server的容器，后面将要使用的Hive会使用这里的MySQL作为metastore数据库存储元数据，但如果仅仅想构建基础Hadoop集群，可省略此步。

## 构建MySQL镜像

通过查询Docker Hub可以找到很多Oracle官方制作好的MySQL镜像，所以这里我们采用版本为5.6.40的镜像作为基础来构建镜像，监听3306端口，并前台运行mysqld守护进程。以下是Dockerfile，写的比较简单：

```docker
FROM mysql:5.6.40

ENV USER=root

# 添加测试用户mysql，密码mysql
RUN echo "mysql:mysql" | chpasswd 

EXPOSE 3306
CMD ["mysqld"]
```

基于官方镜像进行构建省了很多事儿，几乎没有需要调整的配置项，直接使用默认的就可以了。那么直接编写编排文件吧，如下所示：

```
version: "3"
services: 
  mysql:
    image: leesper/mysql
    container_name: mysql
    volumes:
      # 将宿主机的mysql/data挂载到容器中的/var/lib/mysql作为数据存储点
      - ./mysql/data:/var/lib/mysql  
    network_mode: zoo-net
    command: mysqld  # 容器启动时前台运行mysqld
    environment:
      MYSQL_ROOT_PASSWORD: root  # 通过环境变量设置root密码
    ports: 
      - "3306:3306"  # 将容器的3306端口映射到宿主机
```

以`docker-compose -f docker-compose-mysql.yml up -d`启动MySQL容器后，还需要进入容器做一些数据库的配置操作，比如修改密码和配置远程访问：
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

至此，MySQL容器就配置完成了，仍然可以通过`docker ps`命令来查看已启动的容器。