FROM mysql:5.6.40

ENV USER=root

# 添加测试用户mysql，密码mysql，并且将此用户添加到sudoers里
ADD my.cnf /etc/mysql/my.cnf

RUN echo "mysql:mysql" | chpasswd 

EXPOSE 3306
CMD ["mysqld"]