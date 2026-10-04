rm -rf MC-Server-Hosting-Software
set -a
source .env
set +a
go build .
chmod +x MC-Server-Hosting-Software
./MC-Server-Hosting-Software