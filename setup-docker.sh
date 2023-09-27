docker network create erudit-net
docker run -d --name pg --net erudit-net -e POSTGRES_PASSWORD=123456 -p 5433:5432 postgres
docker build -t app .
docker run -d --name server --net erudit-net -p 8888:80 app