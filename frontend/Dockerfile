# for development use
FROM node:14-alpine AS base
WORKDIR /app
COPY ./package.json ./
RUN npm install
RUN chown -R node node_modules

FROM base AS copy 
COPY ./src src
COPY ./public public
COPY ./tsconfig.json ./

FROM copy AS builder 
ARG REACT_APP_API_URL
ENV REACT_APP_API_URL $REACT_APP_API_URL
RUN npm run build

FROM node:14-alpine AS prod
WORKDIR /app
COPY --from=builder /app/build build
COPY --from=builder /app/public public
COPY --from=builder /app/node_modules node_modules
COPY --from=builder /app/package.json package.json
CMD ["npm", "run", "start:prod"]
