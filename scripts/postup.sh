PORT=5995

ip link set dev tun0 mtu 1500

IFC=`route -n | awk '{if($1=="0.0.0.0" && $8 != "tun0"){print $8; exit;}}'`
IP=`ip addr show dev "$ifc" | awk '$1 == "inet" { sub("/.*", "", $2); print $2 }'`

iptables -t nat -F
iptables -t nat -A POSTROUTING -o "$IFC" -j SNAT --to-source "$IP"
iptables -P FORWARD ACCEPT
iptables -A INPUT -p tcp --destination-port `expr $PORT + 1` -j DROP