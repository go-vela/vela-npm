./scripts/start-verdaccio.sh

VERDACCIO_TOKEN=$(curl -s \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X PUT --data '{"name": "testuser", "password": "testpass"}' \
  http://localhost:4873/-/user/org.couchdb.user:testuser | jq '.token')

echo "NPM_TOKEN=$VERDACCIO_TOKEN" > .env-docker
echo "CI=vela" >> .env-docker
echo $VERDACCIO_TOKEN | jq -r > .env
