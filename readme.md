# Notification
This application use for send notification area green,yellow or red virus corona. within get radius on 1km from destination.
### Depedency
    - mongodb
    - redis
    - go version 13
### Running application 
```bash
     go run main.go
```

### Http Clinet
#### Send Device Location
```curl
curl --location --request POST 'http://localhost:8082/device/location' \
--header 'Device-ID: fG7jB1vUGec:APA91bFMRQ68WpKA72PMLa-xSwwmZ222Cr_ZKu1EyvF-dhx1cWWxY5BnfngwTGo9NmsY5HSvu2Onwp4YJ_6tmpW-UFXcnlOU5JjOZiLDINTP66O5JKWpWgvDT63_ADuN1btCwXIa3ABG' \
--header 'Content-Type: application/json' \
--data-raw '{
    "latitude": -6.2490307,
    "longitude": 106.8373179
}
'
```
response: http code 200
```response
{
    "message": "success"
}
```
