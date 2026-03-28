FROM mcr.microsoft.com/playwright:v1.51.0-noble

WORKDIR /checks

COPY package.json ./
RUN npm install

COPY playwright.config.ts ./
COPY e2e/ ./e2e/

ENTRYPOINT ["npx", "playwright", "test"]
