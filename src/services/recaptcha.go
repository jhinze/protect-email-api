package services

import (
	"encoding/json"
	"fmt"
	"hinze.dev/home/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var Recaptcha GoogleRecaptchaInterface

type GoogleRecaptchaInterface interface {
	SiteVerify(response string, remoteIp string) (resp *models.RecaptchaResponse, err error)
}

type GoogleRecaptcha struct {
	Secret string
}

type RecaptchaError struct {
	Problem string
}

func (r RecaptchaError) Error() string {
	return r.Problem
}

func (g *GoogleRecaptcha) SiteVerify(response string, remoteIp string) (resp *models.RecaptchaResponse, err error) {
	recaptchaResponse, recaptchaResponseError := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {g.Secret}, "response": {response}, "remoteip": {remoteIp}})
	if recaptchaResponseError != nil {
		log.Println(recaptchaResponseError)
		return resp, recaptchaResponseError
	}
	if recaptchaResponse.StatusCode != 200 {
		logResponseBody(recaptchaResponse)
		return resp, RecaptchaError{
			Problem: fmt.Sprintf("Response status code was %d", recaptchaResponse.StatusCode),
		}
	}
	resp = &models.RecaptchaResponse{}
	decodeError := json.NewDecoder(recaptchaResponse.Body).Decode(resp)
	if decodeError != nil {
		log.Println(decodeError)
		return resp, RecaptchaError{
			Problem: "unable to decode response",
		}
	}
	return resp, err
}

func logResponseBody(response *http.Response) {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	} else {
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
}
