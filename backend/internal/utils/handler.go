package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

type UtilsController struct{}

func NewUtilsController() *UtilsController {
	return &UtilsController{}
}

// FetchIFSC proxies the request to Razorpay IFSC API
func (pController *UtilsController) FetchIFSC(lCtx *gin.Context) {
	lCode := lCtx.Param("code")
	if len(lCode) != 11 {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid IFSC code length", nil)
		return
	}

	lResp, lErr := http.Get(fmt.Sprintf("https://ifsc.razorpay.com/%s", lCode))
	if lErr != nil || lResp.StatusCode != 200 {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid IFSC code or API error", nil)
		return
	}
	defer lResp.Body.Close()

	var lResult map[string]interface{}
	if lErr := json.NewDecoder(lResp.Body).Decode(&lResult); lErr != nil {
		response.Error(lCtx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to parse API response", lErr.Error())
		return
	}

	response.Success(lCtx, http.StatusOK, "Fetched IFSC details", lResult)
}

// FetchPincode proxies the request to Postal Pincode API
func (pController *UtilsController) FetchPincode(lCtx *gin.Context) {
	lCode := lCtx.Param("code")
	if len(lCode) != 6 {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode length", nil)
		return
	}

	lTr := &http.Transport{
		ForceAttemptHTTP2: false,
		// Disable HTTP/2 by providing an empty map for TLSNextProto
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	lClient := &http.Client{Transport: lTr}

	lResp, lErr := lClient.Get(fmt.Sprintf("https://api.postalpincode.in/pincode/%s", lCode))
	if lErr != nil {
		fmt.Println("Pincode API error:", lErr)
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode or API error", lErr.Error())
		return
	}
	if lResp.StatusCode != 200 {
		fmt.Println("Pincode API status code:", lResp.StatusCode)
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode or API error", fmt.Sprintf("status: %d", lResp.StatusCode))
		return
	}
	defer lResp.Body.Close()

	var lResult []map[string]interface{}
	if lErr := json.NewDecoder(lResp.Body).Decode(&lResult); lErr != nil {
		fmt.Println("Pincode JSON decode error:", lErr)
		response.Error(lCtx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to parse API response", lErr.Error())
		return
	}

	fmt.Println("Pincode parsed result:", lResult)

	if len(lResult) > 0 && lResult[0]["Status"] == "Success" {
		lPostOffices, lOk := lResult[0]["PostOffice"].([]interface{})
		if lOk && len(lPostOffices) > 0 {
			lFirstPO := lPostOffices[0].(map[string]interface{})
			lData := map[string]interface{}{
				"state":    lFirstPO["State"],
				"district": lFirstPO["District"],
				"country":  lFirstPO["Country"],
			}
			response.Success(lCtx, http.StatusOK, "Fetched Pincode details", lData)
			return
		}
	}

	response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode", nil)
}
