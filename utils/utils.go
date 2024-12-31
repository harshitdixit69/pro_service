package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const Salt = "randomUniqueSaltValue!@#$%^&*()" // This should be randomly generated for each user

func HashPasswordWithSalt(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(Salt + password))      // Combine salt and password
	return hex.EncodeToString(hasher.Sum(nil)) // Convert hash bytes to a hex string
}
func MaskAadhaar(aadhaar string) (string, error) {
	// Aadhaar must be 12 digits
	if len(aadhaar) != 12 {
		return "", fmt.Errorf("invalid Aadhaar number length")
	}

	// Ensure the input is numeric
	if _, err := fmt.Sscanf(aadhaar, "%d", new(int)); err != nil {
		return "", fmt.Errorf("Aadhaar number must contain only digits")
	}

	// Replace the middle 6 digits with asterisks
	masked := aadhaar[:4] + strings.Repeat("*", 6) + aadhaar[10:]
	return masked, nil
}

func Encoded(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Decoded(data string) string {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return ""
	}
	return string(decoded)
}
func SendSms(phoneNo, code string) map[string]interface{} {
	url := ""
	method := "POST"
	var to string
	var payload *strings.Reader

	if code == "" {
		to = fmt.Sprint("To=%2B91-", phoneNo, "&Channel=sms")
	} else {
		url = ""
		to = fmt.Sprint("To=%2B91-", phoneNo, "&Code=", code)
	}
	payload = strings.NewReader(to)
	// Twilio credentials
	accountSID := ""
	authToken := ""

	// Create the Basic Auth token
	auth := base64.StdEncoding.EncodeToString([]byte(accountSID + ":" + authToken))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+auth)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)
	return resMap
}
