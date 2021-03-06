package utils

var PROD_BOT_ID = "918359646828388384"
var APP_ENV = getenv("BOT_ENV", "dev")

var DEFAULT_PREFIX = "wump "
var PREFIX = getenv("PREFIX", DEFAULT_PREFIX)

var WUMPUS_GUILD_ID = "918354200482709505"
var GUILD_ID = getenv("GUILD_ID", WUMPUS_GUILD_ID)

var CDN_CHANNEL_ID = getenv("CDN_CHANNEL", "918725182330400788")
var PICS_CHANNEL_ID = getenv("PICS_CHANNEL", "918355152493215764")
var LOGS_CHANNEL_ID = getenv("LOGS_CHANNEL", "918952346975862824")
var VERIFICATION_CHANNEL_ID = getenv("VERIFICATION_CHANNEL", "918932836428419163")

var TEAM_ROLE_ID = getenv("TEAM_ROLE", "918354701337116703")
var OWNER_ROLE_ID = getenv("OWNER_ROLE", "918355466894065685")

var VALID_REACTIONS = []string{"👍", "👎"}

var DEFAULT_POINTS_WORKER_HOST = "https://points.wumpus.club"
var POINTS_WORKER_HOST = getenv("POINTS_WORKER_HOST", DEFAULT_POINTS_WORKER_HOST)
var POINTS_WORKER_SECRET = getenv("POINTS_WORKER_SECRET", "provide_in_env")

var EXT_TO_MIME = map[string]string{"gif": "image/gif", "png": "image/png", "jpg": "image/jpeg", "jpeg": "image/jpeg", "webp": "image/webp"}

var VIDEO_FORMATS = []string{"mp4", "mov", "mkv", "avi", "wmv", "webm", "flv"}

var CDN_ENDPOINT = getenv("CDN_ENDPOINT", "cdn.wumpus.club")
var MINIO_ENDPOINT = getenv("MINIO_ENDPOINT", "provide_in_env")
var MINIO_ACCESS_KEY = getenv("MINIO_ACCESS_KEY", "provide_in_env")
var MINIO_SECRET_KEY = getenv("MINIO_SECRET_KEY", "provide_in_env")
