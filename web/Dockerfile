# Build environment
FROM node:lts-alpine as build-stage

# make the 'app' folder the current working directory
WORKDIR /app

# copy both 'package.json' and 'package-lock.json' (if available)
COPY ./school-front/ ./

# install project dependencies
RUN npm install


# build app for production with minification
RUN npm run build

# Production environment
FROM nginx:1.18.0-alpine as production-stage

RUN mkdir /app

# Copy the conf and content to docker host
COPY --from=build-stage /app/dist /app
COPY ./site.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]