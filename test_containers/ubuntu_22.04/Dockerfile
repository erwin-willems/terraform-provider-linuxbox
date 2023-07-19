FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y openssh-server && \
    apt-get clean && \
    rm --recursive --force /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir /var/run/sshd && \
    echo 'root:root' |chpasswd && \
    sed -r -i 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -r -i 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config && \
    sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd && \
    ssh-keygen -A && \
    mkdir /root/.ssh || true

EXPOSE 22

CMD    ["/usr/sbin/sshd", "-D"]
