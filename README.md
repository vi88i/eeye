# eeye
eagle eye stock screener

# Flow

- User provides a list of stock symbols
- User implements the screening algorithm (pipelining)

# Data Fetch

- Check if data is present in the local cache
  If not:
    Fetch entire historical data
  Else:
    Fetch the missing data since the last available data

Question:
- What about stock split?

# Setup

```cmd
docker pull timescale/timescaledb:latest-pg14

docker run -d --name eeye-db `
  -p 5432:5432 `
  -v eeye-vol:/var/lib/postgresql/data `
  -e POSTGRES_PASSWORD=root `
  -e POSTGRES_USER=admin `
  -e POSTGRES_DB=eeye `
  timescale/timescaledb:latest-pg14

docker exec -it eeye-db psql -U admin -d eeye
```