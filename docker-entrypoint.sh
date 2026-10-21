#!/bin/sh
# Контейнер стартует от root только чтобы привести права на смонтированный
# диск к пользователю app, дальше процесс работает без root.
set -e

if [ "$(id -u)" = "0" ]; then
  mkdir -p "$DATA_DIR"
  chown -R app:app "$DATA_DIR" 2>/dev/null || true
  exec su-exec app /app/portfolio "$@"
fi

exec /app/portfolio "$@"
