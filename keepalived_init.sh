sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/' /etc/sysctl.conf
echo net.ipv4.ip_nonlocal_bind = 1 >> /etc/sysctl.conf
sysctl -p
/bin/sed -i "s/{{VIRTUAL_IP}}/${VIRTUAL_IP}/g" /etc/keepalived/keepalived.conf
/bin/sed -i "s/{{VIRTUAL_MASK}}/${VIRTUAL_MASK}/g" /etc/keepalived/keepalived.conf
/bin/sed -i "s/{{MASTER_IP}}/${MASTER_IP}/g" /etc/keepalived/keepalived.conf
/bin/sed -i "s/{{BACKUP_IP}}/${BACKUP_IP}/g" /etc/keepalived/keepalived.conf
/bin/sed -i "s/{{NF_NAME}}/${NF_NAME}/g" /etc/keepalived/keepalived.conf
/bin/sed -i "s/{{NF_NAME}}/${NF_NAME}/g" /etc/keepalived/notify.sh
/bin/sed -i "s/{{INTERFACE}}/${INTERFACE}/g" /etc/keepalived/keepalived.conf
systemctl enable --now keepalived
service keepalived start