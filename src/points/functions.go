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
	request, _ := httpClient.Get(fmt.Sprintf("%s/balance/%s", utils.POINTS_WORKER_HOST, id))

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
