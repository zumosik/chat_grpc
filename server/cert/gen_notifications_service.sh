rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-notifications-key.pem -out ca-notifications-cert.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=notifications.zumosik.tech"

echo "CA's self-signed certificate"
openssl x509 -in ca-notifications-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-notifications-key.pem -out server-notifications-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=notifications.zumosik.tech"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-notifications-req.pem -days 60 -CA ca-notifications-cert.pem -CAkey ca-notifications-key.pem -CAcreateserial -out server-notifications-cert.pem -extfile server-notifications-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-notifications-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-notifications-key.pem -out client-notifications-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=notifications.zumosik.tech"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-notifications-req.pem -days 60 -CA ca-notifications-cert.pem -CAkey ca-notifications-key.pem -CAcreateserial -out client-notifications-cert.pem -extfile client-notifications-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client-notifications-cert.pem -noout -text