FROM node:22-alpine AS builder

ARG VITE_API_BASE_URL

WORKDIR /frontend

COPY package.json package-lock.json ./
RUN npm install

COPY . ./

ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN npm run build

FROM nginx:alpine

COPY --from=builder /frontend/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 5173



