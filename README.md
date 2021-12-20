# Wumpus Bot

Needs `DISCORD_TOKEN` environment variable

.env file will be auto loaded by godotenv/autoload

Configurable env variables (Value here is the default)

```env
GUILD_ID="918354200482709505"

CDN_CHANNEL="918725182330400788"
PICS_CHANNEL="918355152493215764"
LOGS_CHANNEL="918952346975862824"
VERIFICATION_CHANNEL="918932836428419163"

TEAM_ROLE="918354200482709505"
OWNER_ROLE="918355466894065685"

POINTS_WORKER_HOST="host_for_worker"
POINTS_WORKER_SECRET="provide_in_env"
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

## Running for production

In order to enable all features and commands you will need to set the env to production

`BOT_ENV=production`

## Running in dev

The default envrionment is dev so no need to change that, but set a different prefix with the `PREFIX` env

This will run the bot without the picture logic, or the verification logic, won't send join messages, and won't keep track of WumpCoin.

### Running with different channels and bot

If you want to run the bot with all the features enabled you should change the channel ID's for the verification channel and pics channel to different channels which users cannot see (or users will get mentioned in there when they join)

Then set the env to production, change the prefix and optionally but recommended set the token to a dev bot
