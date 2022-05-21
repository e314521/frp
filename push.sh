#cp ./frps /usr/bin/frps
#mkdir /etc/frp/
#cp ./conf/frps.ini /etc/frp/frps.ini
#cp ./conf/frps.service /etc/systemd/system/


/usr/bin/cp -f./frps /usr/bin/frps 
/usr/bin/cp -f ./frps.ini /etc/frp/frps.ini 
chmod 777 /usr/bin/frps
chmod 777 /etc/frp/frps.ini
setcap cap_net_bind_service=+eip /usr/bin/frps
systemctl enable frps
systemctl restart frps
systemctl status frps


export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./client/libx64
go run ./cmd/frpc/main.go -c ./conf/frpc.ini
git clone git@github.com:e314521/frp.git
git commit -a -m "开启CGO"
git push origin master
