FROM node:22-slim
RUN apt-get update && apt-get install git -y && apt-get clean cache && rm -rf /var/cache/apt /var/lib/apt/lists/*
WORKDIR /
RUN git clone -b main https://github.com/majeinfo/etcd-manager.git
WORKDIR /etcd-manager
RUN npm create vite@latest frontend -- --template react && rm -rf /root/.npm
WORKDIR frontend
RUN cp -f ../frontend-src/App.* .
RUN npm install && rm -rf /root/.npm
RUN npm run build && rm -rf /root/.npm
RUN npm install -g serve && rm -rf /root/.npm

ENV LISTEN_PORT "3000"
ENTRYPOINT ["serve"]
#CMD ["-s", "--listen", "${LISTEN_PORT}", "dist"]
CMD -s --listen ${LISTEN_PORT} dist


