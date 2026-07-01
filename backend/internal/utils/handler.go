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
func (c *UtilsController) FetchIFSC(ctx *gin.Context) {
	code := ctx.Param("code")
	if len(code) != 11 {
		response.Error(ctx, http.StatusBadRequest, "BAD_REQUEST", "Invalid IFSC code length", nil)
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://ifsc.razorpay.com/%s", code))
	if err != nil || resp.StatusCode != 200 {
		response.Error(ctx, http.StatusBadRequest, "BAD_REQUEST", "Invalid IFSC code or API error", nil)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to parse API response", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Fetched IFSC details", result)
}

// FetchPincode proxies the request to Postal Pincode API
func (c *UtilsController) FetchPincode(ctx *gin.Context) {
	code := ctx.Param("code")
	if len(code) != 6 {
		response.Error(ctx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode length", nil)
		return
	}

	tr := &http.Transport{
		ForceAttemptHTTP2: false,
		// Disable HTTP/2 by providing an empty map for TLSNextProto
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(fmt.Sprintf("https://api.postalpincode.in/pincode/%s", code))
	if err != nil {
		fmt.Println("Pincode API error:", err)
		response.Error(ctx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode or API error", err.Error())
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("Pincode API status code:", resp.StatusCode)
		response.Error(ctx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode or API error", fmt.Sprintf("status: %d", resp.StatusCode))
		return
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Pincode JSON decode error:", err)
		response.Error(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to parse API response", err.Error())
		return
	}

	fmt.Println("Pincode parsed result:", result)

	if len(result) > 0 && result[0]["Status"] == "Success" {
		postOffices, ok := result[0]["PostOffice"].([]interface{})
		if ok && len(postOffices) > 0 {
			firstPO := postOffices[0].(map[string]interface{})
			data := map[string]interface{}{
				"state":    firstPO["State"],
				"district": firstPO["District"],
				"country":  firstPO["Country"],
			}
			response.Success(ctx, http.StatusOK, "Fetched Pincode details", data)
			return
		}
	}

	response.Error(ctx, http.StatusBadRequest, "BAD_REQUEST", "Invalid Pincode", nil)
}
