docker compose up -d npm

until npm ping --registry http://localhost:4873; do
  >&2 echo "Verdaccio is unavailable - sleeping"
  sleep 1
done
