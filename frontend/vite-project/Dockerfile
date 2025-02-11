FROM node:latest
WORKDIR /app
COPY package*.json ./
RUN npm install -g npm@latest
COPY . .
EXPOSE 5173
CMD ["npm","run","dev","--","--host","45.141.79.120"]
