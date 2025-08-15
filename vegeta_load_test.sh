#!/bin/bash
# Script de teste de carga usando Vegeta
# Uso: ./vegeta_load_test.sh [URL] [DURACAO] [RATE]
# Exemplo: ./vegeta_load_test.sh http://localhost:8080/rota 30s 100

URL=${1:-http://localhost:8080/}
DUR=${2:-30s}
RATE=${3:-100}

if ! command -v vegeta &> /dev/null; then
  echo "Vegeta não está instalado. Instale com: go install github.com/tsenart/vegeta@latest"
  exit 1
fi

echo "GET $URL" | vegeta attack -duration=$DUR -rate=$RATE | vegeta report
