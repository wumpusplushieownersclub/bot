package points

type PointsAccount struct {
	User   string `json:"user"`
	Points int64  `json:"points"`
}

type Leaderboard []PointsAccount
