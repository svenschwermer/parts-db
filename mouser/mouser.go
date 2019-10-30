package mouser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/svenschwermer/parts-db/config"
)

const searchAPI = "https://api.mouser.com/api/v1/search/partnumber"

// Part holds the part info returned by the Mouser API
type Part struct {
	DataSheetURL           string `json:"DataSheetUrl"`
	Description            string
	ImagePath              string
	Category               string
	Manufacturer           string
	ManufacturerPartNumber string
	MouserPartNumber       string
	ProductDetailURL       string `json:"ProductDetailUrl"`
}

func GetPart(pn string) (*Part, error) {
	fail := func(err error) (*Part, error) {
		log.Println(err)
		return nil, err
	}
	var mouserReq struct {
		SearchByPartRequest struct {
			MouserPartNumber  string `json:"mouserPartNumber"`
			PartSearchOptions string `json:"partSearchOptions"`
		} `json:"SearchByPartRequest"`
	}
	mouserReq.SearchByPartRequest.MouserPartNumber = pn
	mouserReq.SearchByPartRequest.PartSearchOptions = "Exact"

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(mouserReq); err != nil {
		return fail(fmt.Errorf("mouser: request encoding failed: %s", err))
	}
	url := searchAPI + "?apiKey=" + config.Env.MouserAPIKey
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return fail(fmt.Errorf("mouser: request failed: %s", err))
	}
	defer resp.Body.Close()

	var mouserResp struct {
		Errors []struct {
			ID                    int `json:"Id"`
			Code                  string
			Message               string
			ResourceKey           string
			ResourceFormatString  string
			ResourceFormatString2 string
			PropertyName          string
		}
		SearchResults struct {
			Parts []Part
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&mouserResp); err != nil {
		return fail(fmt.Errorf("mouser: response decoding failed: %s", err))
	}

	if len(mouserResp.Errors) > 0 {
		return fail(fmt.Errorf("mouser: error response: %+v", mouserResp.Errors))
	}
	if len(mouserResp.SearchResults.Parts) == 0 {
		return fail(errors.New("part not found"))
	}
	return &mouserResp.SearchResults.Parts[0], nil
}
