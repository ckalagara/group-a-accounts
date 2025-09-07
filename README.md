# group-a-accounts



## dev setup

### docker

```
% docker build -t group-a-accounts .  
```

```
docker run --network app-network --name group-a-accounts -p 8085:8085 group-a-accounts
2025/09/07 19:28:52 Starting service group-a-accounts
2025/09/07 19:28:52 Setting up routes
2025/09/07 19:28:52 Listening on :8085

```

### client

#### Health

```
% curl http://localhost:8085/health
{"status":"OK","type":"HealthCheck","description":"Service is up and running"}   
```