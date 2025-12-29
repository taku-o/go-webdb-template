#!/bin/bash
# CloudBeaver起動スクリプト

# APP_ENV環境変数が未設定の場合はdevelopをデフォルトとする
export APP_ENV=${APP_ENV:-develop}

# Docker ComposeでCloudBeaverを起動
docker-compose -f docker-compose.cloudbeaver.yml up -d
