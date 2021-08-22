# for development use
FROM node:alpine AS base
WORKDIR /app
COPY ./package.json ./
RUN npm install

FROM base AS copy 
COPY ./src src
COPY ./public public

FROM base AS builder 
RUN npm run build

FROM node:alpine AS prod
WORKDIR /app
COPY --from=builder /app/dist dist
COPY --from=builder /app/node_modules node_modules
COPY --from=builder /app/package.json package.json
CMD ["npm", "run", "start:prod"]
