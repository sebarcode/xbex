# Xibar Example Backend

# Prequisite
- Caddy 2 installed
- Db prep executed
- Db migration executed for every db changes

# Preparation
- edit devops/local/app.yml as neccessary 
- db and metadata prep => `go run app/setup/main.go --config=devops/local/app.yml`

# Running on local dev
- `caddy run --config=devops/local/Caddyfile`
- run core => `go run app/core/main.go --config=devops/local/app.yml`
- run inventory => `go run app/invent/main.go --config=devops/local/app.yml`
- run frontend => `cd $front-end-folder; npm run dev`

# Running on docker mode
- create docker for caddy + frontend, core, inventory and pg
- docker compose them in 1 local network
- expose public facing port on docker