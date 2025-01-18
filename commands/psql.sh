docker network create nakama-network || true
docker run -d \
  --name postgres \
  --network nakama-network \
  -e POSTGRES_DB=nakama \
  -e POSTGRES_PASSWORD=localdb \
  -p 5432:5432 \
  postgres:17-alpine

echo "Waiting for the PostgreSQL container to be healthy..."
retries=0
until docker exec postgres pg_isready -U postgres -d nakama || [ $retries -eq 10 ]; do
  echo "PostgreSQL is not ready yet. Retrying..."
  retries=$((retries + 1))
  sleep 3
done
if [ $retries -eq 10 ]; then
  echo "PostgreSQL failed to become healthy"
  exit 1
fi