FROM node:25-trixie-slim

RUN apt-get update \
    && apt-get install --yes --no-install-recommends openjdk-21-jre-headless \
    && rm -rf /var/lib/apt/lists/* \
    && npm install --global firebase-tools@15.32.1

WORKDIR /workspace

COPY firebase.json ./firebase.json

EXPOSE 4000 4400 8080 9150

CMD ["firebase", "emulators:start", "--only", "firestore", "--project", "demo-restaurant-order-server", "--import", "/data", "--export-on-exit", "/data"]
