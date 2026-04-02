FROM node:20-bookworm-slim

WORKDIR /app

ENV UI_HOST=0.0.0.0

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev"]
