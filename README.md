# social-network
[![ci](https://github.com/vvkh/social-network/actions/workflows/ci.yml/badge.svg)](https://github.com/vvkh/social-network/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/vvkh/social-network/branch/main/graph/badge.svg?token=Z1O17LP4FF)](https://codecov.io/gh/vvkh/social-network)

### Setup up dev envrionment
Assuming you have docker-compose and go installed.
```
make install
```

### Start dev environment
This would start local db and run all required migrations.
```
make env
```

### Start local server
Check `.env` file and tweak configs if you need.
#### Run in docker-compose
```
make up-docker
```

#### Compile binaries and run locally
```
make up
```

#### Run/Debug in IDE
Run `cmd/main.go`, environment variables would load automatically.

## Configuration
| env | purpose | example | 
| --- | ------- | ------ |
| SERVER_ADDRESS | listen address as `host:port` | `:80` | 
| TEMPLATES_DIR | path to directory with templates| `./templates` | 
| AUTH_SECRET | password for JWT | `secret` | 
| DB_URL | mysql connection url | `user:password@tcp(localhost:3306)/social_network` |