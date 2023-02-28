# mtls

provides an example on how to enable mTLS between a client and server using a GRPC greeter service

## server

uses the cert watcher to check for new CERTs.

Also ClientAuth is enabled and ClientCAs are provided
On top a middleFunc is used to help validate the client certificate

## client

for long lived sessions the Identity is validated using the latest CERT. if CERT would change the new sessions will use it while exisiting session continue with the old info since the identity was validated when the initial session was established.

As a result the client does not use the certwatcher, When the connection closes it is assumed the new CERT (if changed is picked up)