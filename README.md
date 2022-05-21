# How to run Monolith Program
## install mysql
Once complete the installation, you shoule modify the root's password with cxDgTq9K
## build an run
git clone https://github.com/zchaoyu1126/mini-tiktok.git
cd Monolith
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod init mini-tiktok
go mod tidy
go build .
./mini-tiktok.exe