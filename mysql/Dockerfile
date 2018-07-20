FROM ubuntu:18.04

ENV MYSQL_URL=https://cdn.mysql.com//Downloads/MySQL-5.6/mysql-5.6.40.tar.gz

RUN apt-get update \
    && apt-get install -y --no-install-recommends cmake net-tools curl \
    && rm -rf /var/lib/apt/lists/* \
    && set -x \
    && curl -fSL "$MYSQL_URL" -o /tmp/mysql.tar.gz \
    && tar -xvf /tmp/mysql.tar.gz -C /usr/local \
    && rm -rf /tmp/mysql.tar.gz* \
    && cd /usr/local/mysql-5.6.40 \
    && cmake . -DCMAKE_INSTALL_PREFIX=/usr/local/mysql-5.6.40 -DMYSQL_DATADIR=/usr/local/mysql-5.6.40/data -DSYSCONFDIR=/etc -DWITH_INNOBASE_STORAGE_ENGINE=1 -DWITH_ARCHIVE_STORAGE_ENGINE=1 -DWITH_BLACKHOLE_STORAGE_ENGINE=1 -DWITH_PARTITION_STORAGE_ENGINE=1 -DWITH_PERFSCHEMA_STORAGE_ENGINE=1 -DWITHOUT_EXAMPLE_STORAGE_ENGINE=1 -DWITHOUT_FEDERATED_STORAGE_ENGINE=1 -DDEFAULT_CHARSET=utf8 -DDEFAULT_COLLATION=utf8_general_ci -DWITH_EXTRA_CHARSETS=all -DENABLED_LOCAL_INFILE=1 -DWITH_READLINE=1 -DMYSQL_UNIX_ADDR=/usr/local/mysql-5.6.40/mysql.sock -DMYSQL_TCP_PORT=3306 -DMYSQL_USER=mysql -DCOMPILATION_COMMENT="lq-edition" -DENABLE_DTRACE=0 -DOPTIMIZER_TRACE=1 -DWITH_DEBUG=1 \
    && make && make install

ENV USER=root

# 添加测试用户mysql，密码mysql，并且将此用户添加到sudoers里
ADD my.cnf /etc/my.cnf

RUN useradd mysql \ 
    && echo "mysql:mysql" | chpasswd \
    && echo "mysql   ALL=(ALL)       ALL" >> /etc/sudoers \
    && chown -R mysql:mysql /usr/local/mysql-5.6.40 \
    && chown -R mysql:mysql /etc/my.cnf \ 
    && cd /usr/local/mysql-5.6.40 && ./scripts/mysql_install_db --user=mysql --basedir=/usr/local/mysql-5.6.40 --datadir=/usr/local/mysql-5.6.40/data/

ENV MYSQL_HOME /usr/local/mysql-5.6.40

# 容器需要开放MySQL 3306端口
EXPOSE 3306
