FROM node:16-alpine
ENV REACT_APP_BE_URL=http://localhost:9000
ENV REACT_APP_BE_WS=localhost:9000
WORKDIR /app
COPY . .
RUN npm i
CMD ["npm", "run", "start"]

