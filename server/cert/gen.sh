rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-auth-key.pem -out ca-auth-cert.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=zumosik.tech"

echo "CA's self-signed certificate"
openssl x509 -in ca-auth-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-auth-key.pem -out server-auth-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=zumosik.tech"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-auth-req.pem -days 60 -CA ca-auth-cert.pem -CAkey ca-auth-key.pem -CAcreateserial -out server-auth-cert.pem -extfile server-auth-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-auth-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-auth-key.pem -out client-auth-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=zumosik.tech"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-auth-req.pem -days 60 -CA ca-auth-cert.pem -CAkey ca-auth-key.pem -CAcreateserial -out client-auth-cert.pem -extfile client-auth-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client-auth-cert.pem -noout -text