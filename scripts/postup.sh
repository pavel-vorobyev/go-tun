SERVERPORT=5995
ip link set dev tun0 mtu 1500
ifc=`route -n | awk '{if($1=="0.0.0.0" && $8 != "tun0"){print $8; exit;}}'`
INTERFACE="$ifc"
ip=`ip addr show dev "$ifc" | awk '$1 == "inet" { sub("/.*", "", $2); print $2 }'`
SERVERIP=$ip
iptables -t nat -F
iptables -t nat -A POSTROUTING -o "$ifc" -j SNAT --to-source $ip
iptables -P FORWARD ACCEPT
iptables -A INPUT -p tcp --destination-port `expr $SERVERPORT + 1` -j DROP