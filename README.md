# Mini-tiktok

Hi there, this is a project for practice during the Byte Dance Youth Camp for backend.

This repository is the result of Team BDXA0429.

For Monolith, we choose gin+gorm+mysql+redis+mvc framework to complete the work.

## Monolith

### Project structure

```
Monolith
├── app                   
│   ├── controller        
│   ├── service           
│   └── entity             
|   └── dao
├── common
│   ├── auth              
│   ├── config            
│   ├── db                // mysql & redis init
│   ├── errors            // customed errors
│   ├── logger            
│   └── utils             
├── router                
│   ├── middleware        
│   └── router.go
├── go.mod                
├── go.sum
├── main.go               
├── run.sh
```

### How to run

#### 1. Install mysql

The mysql installation tutorial is omitted here.

Notice: you need replace the `mysql_word` with your own password.

![image-20220527204317266](.\assets\image-20220527204317266.png)



#### 2. Install redis

```shell
wget https://download.redis.io/redis-stable.tar.gz
tar -zvxf redis-stable.tar.gz
mv redis-stable /usr/local/redis
cd /usr/local/redis
make
make test
make PREFIX=/usr/local/redis install
./bin/redis-server&
```

#### 3. Build an run

```shell
git clone https://github.com/zchaoyu1126/mini-tiktok.git
cd Monolith
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod init mini-tiktok
go mod tidy
go build .
sh run.sh
```

### Register a service

```shell
vim /lib/systemd/system/mini-tiktok.service
```

```shell
[Unit]
Description=mini-tiktok

[Service]
Type=simple
Restart=always
RestartSec=3s
ExecStart=sh root/mini-tiktok/Monolith/run.sh

[Install]
WantedBy=multi-user.target
```

After completing the registration, you can use the following command to **start** or **restart** or **stop** the service rather than `sh run.sh`.

```shell
service mini-tiktok start
service mini-tiktok restart
service mini-tiktok stop
service mini-tiktok status
```

### Needs to be done

- [x] use snow flake to generate user_id

- [x] use `sha256(user_id-timestamp)` to generate token 

- [x] use redis to store token.

- [x] use redis to store username to speed up the operation of querying whether the user name exists.

- [ ] set token expire time

- [ ] init redis when the system start.

- [ ] user logout, delete token

- [ ] custom errors

- [ ] add logger

  

