./mkcert -cert-file server-cert.pem -key-file server-key.pem e314521.xyz *.e314521.xyz 192.168.101.220 127.0.0.1

./mkcert -client -cert-file client-cert.pem -key-file client-key.pem e314521.xyz *.e314521.xyz 192.168.101.220 127.0.0.1

./mkcert -pkcs12 -cert-file client.pfx e314521.xyz *.e314521.xyz 192.168.101.220 127.0.0.1

./mkcert -install
