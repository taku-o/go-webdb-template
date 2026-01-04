# 依存関係インストールステージ
FROM node:22-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

# ビルドステージ
FROM node:22-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# 環境変数を設定（ビルド時に必要）
# NEXT_PUBLIC_で始まる環境変数はビルド時にクライアント側のJavaScriptに埋め込まれる
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_API_KEY
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_API_KEY=${NEXT_PUBLIC_API_KEY}

RUN npm run build

# 実行ステージ（開発モード）
FROM node:22-alpine AS dev

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NODE_ENV=development
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL:-http://api:8080}

EXPOSE 3000

CMD ["npm", "run", "dev"]

# 実行ステージ（本番モード）
FROM node:22-alpine AS production

WORKDIR /app

ENV NODE_ENV=production

# ビルド成果物をコピー
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=deps /app/node_modules ./node_modules

# 既存のnodeユーザーを使用（node:22-alpineには既にuid/gid 1000のnodeユーザーが存在）
RUN chown -R node:node /app

USER node

EXPOSE 3000

CMD ["npm", "start"]
