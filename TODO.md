1. determine what data we want from the API DONE

7 day high low
--------------
id
lat
long
date
high
low


2. define db table in .sql file DONE



3. create interface to write/read data to table
4. create rest package to call service; service will callweather API 
5. implement http api server
6. work on html/template



ENTER CITY to get LAT/LONG - examples
https://geocoding-api.open-meteo.com/v1/search?name=Berlin&count=10&language=en&format=json

FORECAST
https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current=temperature_2m&daily=temperature_2m_max,temperature_2m_min&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch&timezone=America%2FNew_York&forecast_days=16
