FROM mysql:5.6.40

ENV USER=root

# 添加测试用户mysql，密码mysql
RUN echo "mysql:mysql" | chpasswd 

EXPOSE 3306
CMD ["mysqld"]