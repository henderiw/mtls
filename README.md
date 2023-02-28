# mtls

provides an example on how to enable mTLS between a client and server using a GRPC greeter service

## use

## standalone

start the server

```
go run *.go server --local
```

start the client

```
go run *.go client --local
```

### kubernetes

install cert manager

```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
```

create a proxy namespace

```
kubectl create ns proxy
```

deploy the cert/issuer infrastructure

```
kubectl apply -f config/certs/single-subca
```

deploy the server app

```
kubectl apply -f config/server
```

deploy the client app

```
kubectl apply -f config/client
```

## server

uses the cert watcher to check for new CERTs.

Also ClientAuth is enabled and ClientCAs are provided
On top a middleFunc is used to help validate the client certificate

## client

for long lived sessions the Identity is validated using the latest CERT. if CERT would change the new sessions will use it while exisiting session continue with the old info since the identity was validated when the initial session was established.

As a result the client does not use the certwatcher, When the connection closes it is assumed the new CERT (if changed is picked up)