# Dockerized version of Temporal

Benefits:
- Support for ARM processors
- It contains `reset` endpoint - `POST http://127.0.0.1:1323/reset`.

## How use it

In your docker-compose.yaml file add new service:

```yaml
version: '3'

services:
  temporal-test-server:
    image: ghcr.io/cv65kr/temporal-test-server:latest
    platform: linux/amd64
```
