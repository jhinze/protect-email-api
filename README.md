# Protect Email API
A microservice that returns an email address upon the successful validation of a reCAPTCHA response.   

In order to prevent non-humans from seeing a personal email address.

## Usage
_request_   
```  
GET /v1/email?token=<reCAPTCHA token>   
```     

_respose_
```
{ 
  "email":"foo@bar.baz"
}
```

see it in action at www.hinze.dev

## Required Environment Variables
`RECAPTCHA_SECRET=<reCAPTCHA secret>`   
`PROTECTED_EMAIL=<email to return>`

## Docker
`docker run -p <port>:80 -e RECAPTCHA_SECRET=<reCAPTCHA secret> -e PROTECTED_EMAIL=<email to return> jhinze/protect-email-api`
