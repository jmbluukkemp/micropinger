# MicroPinger 1.0

## What?
Simple microservice (6.28MB docker image) which can act as a server or client for monitoring http endpoints.  
Client component will post to a Slack channel if the required endpoint is unreachable.  
Server component will echo "pong" on a HTTP request.  

## How?

### Server
To start in server mode, please set environment variable "mode" to "server"  
The server start running on port :8000 and accepts any GET request

### Client
To start in client mode, please set environment variable "mode" to "client"  
Set environment variable "slack" to a Slack webhook endpoint  
Eg: 
```bash
export mode=client
export slack=https://hooks.slack.com/services/XXXXXXXX/XXXXXX/XXXXXXXXXXXX
```
This can also be set in a kubernetes secret file (in Base64)  
Example: see client-secret.yaml  

#### Generate base64 string for use in kubernetes secret file
```bash
echo -n "client" | base64 -w0
echo -n "https://hooks.slack.com/services/XXXXXXXX/XXXXXX/XXXXXXXXXXXX" | base64 -w0
```

Create an endpoints.json file containing the endpoints you wish to monitor.  
See endpoints-example.json.  
"id" and "secret" are the x-ibm-client-id and x-ibm-client-secret header values  

### Generate secret from endpoints.json
To generate Kubernetes secret:
```bash
kubectl create secret generic micropinger-endpoints --from-file=endpoints.json
```
