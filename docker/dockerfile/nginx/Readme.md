docker build . -t nginx:with_dockerfile
docker run -d -p 80:80 --name nginx -e ENV=DEV nginx:with_dockerfile
docker logs nginx


