# Dining cryptographers

This is implementation of [dining cryptographers protocol](https://en.wikipedia.org/wiki/Dining_cryptographers_problem).

This app lets you run 3 instances of cryptographers and one instance of server. 
All instances communicates over TLS. 

Below there is a list of ports, that instances are listening on.

| Instance        | Port          | 
| --------------- |:-------------:|
| server          |  443          | 
| cryptographer 1 | 8081          | 
| cryptographer 2 | 8082          | 
| cryptographer 3 | 8083          | 

# Before running this app

<span style="color:red">**IMPORTANT: To run this app you have to have valid go installation.**</span>

To run demonstration first it necessary to generate some certificates. You can do it by hand following instructions below,
but there is also a script that will make it for you. 
1. First of all we need CA root. To generate it we need a key for RSA.
    ```console
    foo@bar:~$ openssl genrsa -out ./resources/CA/rootCA.key 4096
    ```
    Once you have RSA key, you can generate root certificate
    ```console
    foo@bar:~$ openssl req -x509 -new -nodes -key ./resources/CA/rootCA.key -sha256 -days 1024 -out ./resources/CA/rootCA.crt
    ```
    __Important:__ set `Common Name` to `localhost` 
    
2. Now we have our self-signed root certificate. We can generate certificates for server and clients.
    Before we do certificates, we have to generate RSA keys for all of them. Lets begin with the server:
    ```console
    foo@bar:~$ openssl genrsa -out ./resources/server/server.key 4096
    ```
    And then certificate signing request: 
    ```console
    foo@bar:~$ openssl req -new -key ./resources/server/server.key -out ./resources/server/server.csr
    ```
    And finally - certificate:
    ```console
    foo@bar:~$ openssl x509 -req -in ./resources/server/server.csr -CA ./resources/CA/rootCA.crt -CAkey ./resources/CA/rootCA.key -CAcreateserial -out ./resources/server/server.crt -days 500 -sha256
    ```
    
3. Almost the same procedure will be for generating certificates for clients. Below there is an example for client.
    To generate certificates for client 2 and 3 you have to change `1` for `2` and `3` respectively.
    ```console
    foo@bar:~$ openssl openssl genrsa -out ./resources/clients/keys/client1.key 4096
    foo@bar:~$ openssl req -new -key ./resources/clients/keys/client1.key -out ./resources/clients/crts/client1.csr
    foo@bar:~$ openssl x509 -req -in ./resources/clients/crts/client1.csr -CA ./resources/CA/rootCA.crt -CAkey ./resources/CA/rootCA.key -CAcreateserial -out ./resources/clients/crts/client1.crt -days 500 -sha256
    ```

You can do it the same using following script: 
```console
foo@bar:~$ sh ./resources/certificate_generator.sh
```
Now when everything is set up, you can finally run the demo.

# How to run

1. Compiling sources:
    ```console
    foo@bar:~$ go install ./...
    ```
2. To run the demo you need 4 terminal windows. Then run following commands( each one in separate terminal window) in this order.
The order is important, because cryptographer with id 1 initiates flow. 
    ```console
    foo@bar1:~$ sudo runserver 
    foo@bar2:~$ runcryptographer -id=3 
    foo@bar2:~$ runcryptographer -id=2 
    foo@bar2:~$ runcryptographer -id=1 -payed=true 
    ```
    
After running last command, you should see who payed on every terminal. 