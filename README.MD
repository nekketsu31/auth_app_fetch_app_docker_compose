# Efishery Test


## Build

```
$ docker-compose -f ./fetch_app/docker-compose.yml up
```
## Stop

```
$ docker-compose -f ./fetch_app/docker-compose.yml down
```

Open <http://localhost:3000> for auth app.
```
##Register User
curl --location --request POST 'http://localhost:3000/api/register' \
--header 'Content-Type: application/json' \
--data-raw '{
	"phone":"081111111111",
	"name":"user",
	"role":"user"
}'

##Get Detail User
curl --location --request GET 'http://localhost:3000/api/users' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJpYXQiOjE2MjUyMzM0OTIsImV4cCI6MTYyNTI0MzQ5Mn0.Rv3oxTl6l9N10q2iQ-QzHyDDuOrHllmP4zaNm1mDtus'

##Login
curl --location --request POST 'http://127.0.0.1:3000/api/login' \
--header 'Content-Type: application/json' \
--data-raw '{
	"phone":"081111111111",
	"password":"BSWe"
}'
```

Open <http://127.0.0.1:8080/> for fetch app.
```
##Get List Commodity
curl --location --request GET 'http://127.0.0.1:8080/commodity_list' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJpYXQiOjE2MjUyMzU5OTcsImV4cCI6MTYyNTI0NTk5N30.qnve91KgpTyzHw28dadQweEWQioSABA0xlw_Btq-Jhg'

##Commodity Aggragation
curl --location --request GET 'http://localhost:8080/commodity_aggregation' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJpYXQiOjE2MTUxMDExODksImV4cCI6MTYxNTExMTE4OX0.RW9AAFI3S3EDpNYiXq1Ged20ZXF2jzRxl7L3lotHF3A'

##Private Claim
curl --location --request GET 'http://localhost:10000/private_claim' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJpYXQiOjE2MTUxMDExODksImV4cCI6MTYxNTExMTE4OX0.RW9AAFI3S3EDpNYiXq1Ged20ZXF2jzRxl7L3lotHF3A'
```
