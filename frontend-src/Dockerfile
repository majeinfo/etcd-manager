FROM node:22-slim as builder
RUN apt-get update && apt-get install git -y && apt-get clean cache && rm -rf /var/cache/apt /var/lib/apt/lists/*
WORKDIR /
RUN git clone -b main https://github.com/majeinfo/etcd-manager.git
WORKDIR /etcd-manager
RUN npm create vite@latest frontend -- --template react && rm -rf /root/.npm
WORKDIR frontend
RUN cp -f ../frontend-src/App.* ../frontend-src/index.html ./src
RUN npm install && rm -rf /root/.npm
RUN npm run build && rm -rf /root/.npm
RUN cp -f ../frontend-src/runtime-config.js ./dist

FROM node:22-slim
WORKDIR /etcd-manager
COPY --from=builder /etcd-manager/frontend/dist .

ENV LISTEN_PORT "8080"
ENV VITE_API_URL "http://backend:3000"

RUN npm install -g serve && rm -rf /root/.npm

ENTRYPOINT ["serve"]
#CMD ["-s", "--listen", "${LISTEN_PORT}", "dist"]
CMD -s --listen ${LISTEN_PORT} dist


