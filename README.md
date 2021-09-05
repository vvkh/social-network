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

### Backend application
| env | purpose | example | 
| --- | ------- | ------ |
| SERVER_ADDRESS | listen address as `host:port` | `:80` | 
| TEMPLATES_DIR | path to directory with templates| `./templates` | 
| AUTH_SECRET | password for JWT | `secret` | 
| DB_URL | mysql connection url | `user:password@tcp(localhost:3306)/social_network` |
| BCRYPT_COST | cost of hasing used in bcrypt | `10` |


### Local dev environment
#### MySQL main database
| env | purpose | example | 
| --- | ------- | ------ |
| MYSQL_ALLOW_EMPTY_PASSWORD | can empty password be used for root user | `yes` | 
| MYSQL_DATABASE | application db name | `social_network` | 
| MYSQL_USER | application db user | `user` | 
| MYSQL_PASSWORD | application db password | `password` |

#### MySQL replica database
| env | purpose | example | 
| --- | ------- | ------ |
| MYSQL_REPLICA_MAIN_HOST | main db instance host | `db` | 
| MYSQL_REPLICA_MAIN_PORT | main db instance host port | `3306` | 
| MYSQL_REPLICA_USER | replica db user | `user` | 
| MYSQL_REPLICA_PASSWORD | replica db password | `password` |

#### ProxySQL
| env | purpose | example | 
| --- | ------- | ------ |
| PROXY_MYSQL_HOST | main db instance host | `db` | 
| PROXY_MYSQL_PORT | main db instance host port | `3306` | 
| PROXY_REPLICA_MYSQL_HOST | replica db instance host | `db` | 
| PROXY_REPLICA_MYSQL_PORT | replica db instance host port | `3306` | 
| PROXY_SEARCH_PROFILES_HOST | hostgroup id which `/profiles` read request must be routed to, 0 is main, 1 is repilca | `0` | 
| PROXY_HOST | proxysql host | `db_proxy` | 
| PROXY_PORT | proxysql admin port | `6032` | 
| PROXY_USER_PORT | proxysql port to be used by application| `6033` | 
| PROXY_ADMIN_USER | proxysql db user | `user` | 
| PROXY_ADMIN_PASSWORD | proxysql db password | `password` |
