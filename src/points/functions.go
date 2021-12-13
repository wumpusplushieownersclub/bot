package points

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"wumpus/src/utils"
)

func GetPointsAccount(id string) (*PointsAccount, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, reqErr := httpClient.Get(fmt.Sprintf("%s/balance/%s", utils.POINTS_WORKER_HOST, id))

	if reqErr != nil {
		fmt.Println("Error getting balance", reqErr)
		return nil, reqErr
	}

	account := &PointsAccount{}

	body, readErr := ioutil.ReadAll(request.Body)
	if readErr != nil {
		fmt.Println(readErr)
		fmt.Println("Error reading account body", readErr)
		return nil, readErr
	}

	err := json.Unmarshal(body, account)
	if err != nil {
		fmt.Println("Error unmarshalling account", err)
		return nil, err
	}

	return account, nil
}

func GetLeaderboard() (Leaderboard, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, reqErr := httpClient.Get(fmt.Sprintf("%s/leaderboard", utils.POINTS_WORKER_HOST))

	if reqErr != nil {
		fmt.Println("Error getting leaderboard", reqErr)
		return nil, reqErr
	}

	board := Leaderboard{}

	body, readErr := ioutil.ReadAll(request.Body)
	if readErr != nil {
		fmt.Println("Error reading leaderboard body", readErr)
		return nil, readErr
	}

	err := json.Unmarshal(body, &board)
	if err != nil {
		fmt.Println("Error unmarshalling leaderboard", err)
		return nil, err
	}

	return board, nil
}
