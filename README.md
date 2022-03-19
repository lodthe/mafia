# mafia

[Mafia](https://en.wikipedia.org/wiki/Mafia_(party_game)) CLI

## Startup 

You can build binaries by yourself our use built docker images.

### Make

```bash
# git clone git@github.com:lodthe/mafia.git && cd mafia
make all

# Run server 
./bin/server

# Run client
./bin/client
```

### Docker
```bash
# Run server
docker run -d -it -p 80:9000 lodthe/mafia-server

# Run client
docker run -it lodthe/mafia-client --address <server address>
```
