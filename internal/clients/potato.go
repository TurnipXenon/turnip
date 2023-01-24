package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TurnipXenon/potato_api/rpc/potato"

	"github.com/TurnipXenon/turnip/internal/models"
	"github.com/TurnipXenon/turnip/internal/util"
)

type Potato interface {
	RevalidateStaticPath(path string) error
}

type potatoImpl struct {
	token   string
	baseUrl string // todo: set somewhere
}

func NewPotato(flags *models.RunFlags) Potato {
	return &potatoImpl{
		token:   flags.PotatoToken,
		baseUrl: flags.PotatoUrl,
	}
}

func (p *potatoImpl) RevalidateStaticPath(path string) error {
	// todo: rename to revalidate path to make it consistent with potato_api
	// todo: make less hardcoded
	// from https://www.digitalocean.com/community/tutorials/how-to-make-http-requests-in-go
	reqBody := potato.RevalidatePathRequest{
		Path: path,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(nil)
	}
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := fmt.Sprintf("%s/revalidate", p.baseUrl)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", p.token))

	// todo: assess DefaultClient
	// todo: put token on request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}

	if resp.StatusCode != 200 {
		resStr := ""
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
		} else {
			resStr = string(resBody)
		}

		// todo: don't log if it's a 400 error?
		newErr := fmt.Errorf("failed to revalidate with code %s and message %s", resp.Status, resStr)
		util.LogDetailedError(newErr)
		return util.WrapErrorWithDetails(err)
	}

	return nil
}
