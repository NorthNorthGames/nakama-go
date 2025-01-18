docker run -d \
  --name nakama \
  --network nakama-network \
  --link postgres:db \
  --restart always \
  -v config/nakama:/nakama/data \
  -p 7349:7349 \
  -p 7350:7350 \
  -p 7351:7351 \
  registry.heroiclabs.com/heroiclabs/nakama:latest \
  sh -ecx "
    /nakama/nakama migrate up --database.address postgres:localdb@postgres:5432/nakama &&
    exec /nakama/nakama --name nakama1 --database.address postgres:localdb@postgres:5432/nakama --logger.level DEBUG --session.token_expiry_sec 7200
  "

echo "Waiting for the Nakama container to be healthy..."
retries=0
until docker exec nakama /nakama/nakama healthcheck || [ $retries -eq 5 ]; do
  echo "Nakama is not ready yet. Retrying..."
  retries=$((retries + 1))
  sleep 10
done
if [ $retries -eq 5 ]; then
  echo "Nakama failed to become healthy"
  exit 1
fi