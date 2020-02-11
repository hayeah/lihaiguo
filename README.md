判断 IP 是否为大中华局域网, 返回不同 API 节点。

# build & deploy

Build only:

```
make linux
```

Build and deploy:

```
export HOST=howard@34.87.64.224
make deploy
```

Run it:

```
./lihaiguo.bin nodes.yaml
```

# Config

See: [nodes.example.yaml](nodes.example.yaml)

# API

If from a Chinese IP:

```
export API_HOST=34.87.64.224:8888
curl $API_HOST/api_nodes

{"country":"cn","nodes":["123.123.123.1:4442","123.123.123.2:4442"]}
```

If from a US IP:

```
export API_HOST=34.87.64.224:8888
curl $API_HOST/api_nodes

{"country":"us","nodes":["223.123.123.1:4442","223.123.123.2:4442"]}
```

# Maxmind DB

Sign up and download the GeoLite2-Country database:

https://www.maxmind.com/en/geolite2/signup
