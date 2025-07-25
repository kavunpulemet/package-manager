FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y openssh-server sudo && \
    mkdir /var/run/sshd && \
    useradd -m -s /bin/bash user && echo 'user:qwerty' | chpasswd && \
    echo 'user ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers && \
    mkdir -p /home/user/packages && \
    chown user:user /home/user/packages && \
    chmod 755 /home/user/packages

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D"]
