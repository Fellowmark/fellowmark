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

FROM copy AS builder 
RUN npm run build

FROM node:alpine AS prod
WORKDIR /app
COPY --from=builder /app/build build
COPY --from=builder /app/public public
COPY --from=builder /app/node_modules node_modules
COPY --from=builder /app/package.json package.json
CMD ["npm", "run", "start:prod"]
