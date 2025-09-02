#!/bin/sh
set -e

FRONTEND_DIR="/usr/share/nginx/html"
VARS="BACKEND_BASE_URL"

for VAR in $VARS; do
  VALUE=$(printenv "$VAR")
  if [ -n "$VALUE" ]; then
    echo "🔄 Replacing __${VAR}__ with $VALUE"
    grep -rl "__${VAR}__" "$FRONTEND_DIR" | while read -r FILE; do
      sed -i "s|__${VAR}__|$VALUE|g" "$FILE"
    done
  else
    echo "⚠️ Variável $VAR não definida, pulando..."
  fi
done

echo "✅ Variáveis substituídas. Iniciando nginx..."
exec "$@"