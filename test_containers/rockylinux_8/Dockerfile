FROM rockylinux:8

RUN  \
    # Update to latest and install ssh and init scripts.
    # Note: we'll keep the yum cache around to speed up tests.
    yum update --assumeyes && \
    yum install --assumeyes openssh-server initscripts && \
    # Enable systemd:
    rm -f /lib/systemd/system/multi-user.target.wants/*;\
    rm -f /etc/systemd/system/*.wants/*;\
    rm -f /lib/systemd/system/local-fs.target.wants/*; \
    rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
    rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
    rm -f /lib/systemd/system/basic.target.wants/*;\
    rm -f /lib/systemd/system/anaconda.target.wants/*;\
    # Configure SSH:
    mkdir /var/run/sshd && \
    echo 'root:root' |chpasswd && \
    sed -r -i 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -r -i 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config && \
    sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd && \
    ssh-keygen -A && \
    mkdir /root/.ssh || true && \
    systemctl enable sshd

EXPOSE 22

CMD ["/usr/lib/systemd/systemd"]
# sed -i 's/arch=\$basearch/arch=\$basearch\&cc=nl/g' *.repo