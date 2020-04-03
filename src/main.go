package main

import (
	"hinze.dev/home/config"
	"hinze.dev/home/routes"
	"hinze.dev/home/services"
	"log"
	"os"
)

func main() {
	recaptchaSecret, hasSecret := os.LookupEnv("RECAPTCHA_SECRET")
	if !hasSecret || len(recaptchaSecret) == 0 {
		log.Fatalln("Environment variable RECAPTCHA_SECRET is missing or empty")
	}
	services.Recaptcha = &services.GoogleRecaptcha{Secret: recaptchaSecret}

	email, hasEmail := os.LookupEnv("PROTECTED_EMAIL")
	if !hasEmail || len(email) == 0 {
		log.Fatalln("Environment variable PROTECTED_EMAIL is missing or empty")
	}
	config.Email = email

	r := routes.SetupRouter()
	runErr := r.Run()
	if runErr != nil {
		log.Fatal(runErr)
	}
}
