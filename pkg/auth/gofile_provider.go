package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/KhoalaS/godel/pkg/types"
)

const WT string = "4fd6sg89d7s6"

func GofileAuthprovider() (types.Credentials, error) {

	var creds types.Credentials

	accountsUrl := "https://api.gofile.io/accounts"

	referer := "https://gofile.io/"
	origin := "https://gofile.io"

	request, err := http.NewRequest(http.MethodPost, accountsUrl, nil)
	request.Header.Add("Origin", origin)
	request.Header.Add("Referer", referer)
	request.Header.Add("User-Agent", UserAgent)

	if err != nil {
		return creds, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return creds, err
	}

	if response.StatusCode != http.StatusOK {
		return creds, errors.New("could not get accounts data for gofile authprovider")
	}

	var accountResponse AccountResponse

	accountData, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return creds, err
	}

	err = json.Unmarshal(accountData, &accountResponse)
	if err != nil {
		return creds, err
	}

	if accountResponse.Status != "ok" {
		return creds, fmt.Errorf("GofileAuthprovider: got invalid status in accounts response %s", accountResponse.Status)
	}

	creds = types.Credentials{
		Token: accountResponse.Data.Token,
		Headers: map[string]string{
			"authorization": fmt.Sprintf("Bearer %s", accountResponse.Data.Token),
			"origin":        origin,
			"referer":       referer,
			"user-agent":    UserAgent,
		},
	}

	return creds, nil
}

type AccountResponse struct {
	Status string  `json:"status"`
	Data   Account `json:"data"`
}

type Account struct {
	ID         string `json:"id"`
	RootFolder string `json:"rootFolder"`
	Tier       string `json:"tier"`
	Token      string `json:"token"`
}
