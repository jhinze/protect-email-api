# Protect Email API
A microservice that returns an email address upon the successful validation of a reCAPTCHA response.   

In order to prevent non-humans from seeing a personal email address.

## Required Environment Variables
`RECAPTCHA_SECRET=<reCAPTCHA secret>`   
`PROTECTED_EMAIL=<email to returnn>`

## Docker Example
`docker run -p <port>:80 -e RECAPTCHA_SECRET=<reCAPTCHA secret> -e PROTECTED_EMAIL=<email to returnn> jhinze/protect-email-api`
