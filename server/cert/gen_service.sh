#!/bin/bash

SERVICE_NAME=$1


# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-${SERVICE_NAME}-key.pem -out ca-${SERVICE_NAME}-cert.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=${SERVICE_NAME}.zumosik.tech"

echo "CA's self-signed certificate"
openssl x509 -in ca-${SERVICE_NAME}-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-${SERVICE_NAME}-key.pem -out server-${SERVICE_NAME}-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=${SERVICE_NAME}.zumosik.tech"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-${SERVICE_NAME}-req.pem -days 60 -CA ca-${SERVICE_NAME}-cert.pem -CAkey ca-${SERVICE_NAME}-key.pem -CAcreateserial -out server-${SERVICE_NAME}-cert.pem -extfile server-${SERVICE_NAME}-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-${SERVICE_NAME}-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-${SERVICE_NAME}-key.pem -out client-${SERVICE_NAME}-req.pem -subj "/C=US/ST=California/L=San Francisco/O=Company/OU=Some Department/CN=${SERVICE_NAME}.zumosik.tech"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-${SERVICE_NAME}-req.pem -days 60 -CA ca-${SERVICE_NAME}-cert.pem -CAkey ca-${SERVICE_NAME}-key.pem -CAcreateserial -out client-${SERVICE_NAME}-cert.pem -extfile client-${SERVICE_NAME}-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client-${SERVICE_NAME}-cert.pem -noout -text