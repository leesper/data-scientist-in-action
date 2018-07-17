FROM ubuntu:18.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends openjdk-8-jdk openssh-server net-tools curl \
    && rm -rf /var/lib/apt/lists/* \
    && set -x \
    && ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa \
    && cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys \
    && chmod 0600 ~/.ssh/authorized_keys \
    && mkdir -p /var/run/sshd

ENV USER=root
ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64
ENV PATH=${JAVA_HOME}/bin:$PATH

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]