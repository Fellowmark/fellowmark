# for development use
FROM node:alpine AS base
WORKDIR /app
COPY ./package.json ./
RUN npm install
RUN chown -R node node_modules

FROM base AS copy 
COPY ./src src
COPY ./public public
COPY ./tsconfig.json ./

FROM base AS builder 
RUN npm run build

FROM node:alpine AS prod
WORKDIR /app
COPY --from=builder /app/dist dist
COPY --from=builder /app/node_modules node_modules
COPY --from=builder /app/package.json package.json
CMD ["npm", "run", "start:prod"]
