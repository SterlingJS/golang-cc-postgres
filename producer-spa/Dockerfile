FROM node:14-alpine 

ENV PRODUCER_BACKEND https://keda-producer.fmz-c-x-app-aks-01.corp.fmglobal.com/enqueue

WORKDIR /app 

COPY package.json ./

RUN npm install 

COPY . .

EXPOSE 3000 

CMD ["npm", "start"]