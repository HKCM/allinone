FROM nginx
USER root
WORKDIR /bin
ADD index.html /usr/share/nginx/html/index.html
ADD entrypoint.sh /
EXPOSE 80
ENV EnVariable=123
ENTRYPOINT ["/entrypoint.sh"]


