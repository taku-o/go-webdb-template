#!/bin/bash
# Metabase起動スクリプト

# APP_ENV環境変数が未設定の場合はdevelopをデフォルトとする
export APP_ENV=${APP_ENV:-develop}

# Docker ComposeでMetabaseを起動
docker-compose -f docker-compose.metabase.yml up -d
