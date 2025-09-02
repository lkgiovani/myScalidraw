FROM node:22-alpine AS builder
WORKDIR /frontend
COPY package.json package-lock.json ./
COPY .env.production .env
RUN npm install
COPY . ./
RUN npm run build


FROM nginx:alpine
COPY --from=builder /frontend/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY docker-entry.sh /docker-entry.sh
RUN chmod +x /docker-entry.sh
EXPOSE 5173
ENTRYPOINT ["/docker-entry.sh"]
CMD ["nginx", "-g", "daemon off;"]
