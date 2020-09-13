FROM nginx
COPY ./site.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
RUN apt-get update && apt-get install -y \
        curl \
        wget \
        git \
        vim \
        nodejs \
        npm
WORKDIR /app
CMD ["nginx", "-g", "daemon off;"]