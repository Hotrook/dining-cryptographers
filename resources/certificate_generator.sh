#!/usr/bin/env bash

openssl genrsa -out ./resources/CA/rootCA.key 4096
openssl req -x509 -new -nodes -key ./resources/CA/rootCA.key -sha256 -days 1024 -out ./resources/CA/rootCA.crt -config ./resources/conf/csr.conf

openssl genrsa -out ./resources/server/server.key 4096
openssl req -new -key ./resources/server/server.key -out ./resources/server/server.csr -config ./resources/conf/csr.conf
openssl x509 -req -in ./resources/server/server.csr -CA ./resources/CA/rootCA.crt -CAkey ./resources/CA/rootCA.key -CAcreateserial -out ./resources/server/server.crt -days 500 -sha256

for ID in 1 2 3
do
    openssl genrsa -out ./resources/clients/keys/client${ID}.key 4096
    openssl req -new -key ./resources/clients/keys/client${ID}.key -out ./resources/clients/crts/client${ID}.csr -config ./resources/conf/csr.conf
    openssl x509 -req -in ./resources/clients/crts/client${ID}.csr -CA ./resources/CA/rootCA.crt -CAkey ./resources/CA/rootCA.key -CAcreateserial -out ./resources/clients/crts/client${ID}.crt -days 500 -sha256
done

find . -name *.csr -delete