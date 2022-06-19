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
arm-linux-gnueabihf-gcc -fPIC -shared SynReader.c -o libzkident.so

grep -rn "hello,world!" *
scp ./bin/frpc root@192.168.101.220:/oem/
scp ./conf/frpc.ini root@192.168.101.220:/oem/
scp ./serve.pem root@192.168.101.220:/oem/
scp ./serve.key root@192.168.101.220:/oem/

./mkcert -cert-file serve.pem -key-file serve.key -p12-file serve.p12 e314521.xyz *.e314521.xyz 192.168.101.220

./mkcert -client -cert-file client-csr.pem -key-file client-key.pem

./mkcert -install
