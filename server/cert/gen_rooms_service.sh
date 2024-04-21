rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-rooms-key.pem -out ca-rooms-cert.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=rooms.zumosik.tech"

echo "CA's self-signed certificate"
openssl x509 -in ca-rooms-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-rooms-key.pem -out server-rooms-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=rooms.zumosik.tech"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-rooms-req.pem -days 60 -CA ca-rooms-cert.pem -CAkey ca-rooms-key.pem -CAcreateserial -out server-rooms-cert.pem -extfile server-rooms-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-rooms-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-rooms-key.pem -out client-rooms-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=rooms.zumosik.tech"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-rooms-req.pem -days 60 -CA ca-rooms-cert.pem -CAkey ca-rooms-key.pem -CAcreateserial -out client-rooms-cert.pem -extfile client-rooms-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client-rooms-cert.pem -noout -text