package authmiddleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

// GetRandomBase64String create a random byte array and return it base64 encoded
func GetRandomBase64String(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// GetRandomString create a random string with length n
func GetRandomString(n int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	s := make([]byte, n)
	for i := range s {
		s[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(s)
}

// JSONToMap create a map out of JSON string
func JSONToMap(jsonRaw string) (map[string]interface{}, error) {
	flatMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonRaw), &flatMap)
	if err != nil {
		return flatMap, err
	}
	return flatMap, nil
}

// GetIPAdress gets the IP address of requester
// For debug purpose
// Source: https://golangbyexample.com/golang-ip-address-http-request/
func GetIPAdress(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", errors.New("No valid ip found")
}
