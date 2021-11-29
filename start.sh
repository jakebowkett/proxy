sudo setcap CAP_NET_BIND_SERVICE=+eip ./proxy
nohup ./proxy > out.log 2>&1 &
echo $! > pid.txt