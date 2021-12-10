# Wumpus Bot

Needs `DISCORD_TOKEN` environment variable

.env file will be auto loaded by godotenv/autoload

Configurable env variables (Value here is the default)

```env
CDN_CHANNEL="918725182330400788"
PICS_CHANNEL="918355152493215764"
TEAM_ROLE="918354200482709505"
```

## Starting the bot

Run `go run ./src` (Using ./src instead of main.go since there is multiple files in a folder)

If you want reload-on-save you can use [reflex](https://github.com/cespare/reflex)

```sh
go install github.com/cespare/reflex@latest
```

Then run the bot using

```sh
reflex -r '\.go' -s -- sh -c "go run ./src"
```
