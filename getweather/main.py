import openmeteo_requests
import pandas as pd
import requests_cache
from retry_requests import retry
import os
import requests
from datetime import datetime
from dotenv import load_dotenv
from zoneinfo import ZoneInfo
from pathlib import Path
import sys
import time

while True: 
  # Setup the Open-Meteo API client with cache and retry on error
  cache_session = requests_cache.CachedSession('.cache', expire_after = 3600)
  retry_session = retry(cache_session, retries = 5, backoff_factor = 0.2)
  openmeteo = openmeteo_requests.Client(session = retry_session)
  
  # Make sure all required weather variables are listed here
  # The order of variables in hourly or daily is important to assign them correctly below
  url = "https://api.open-meteo.com/v1/forecast"
  params = {
  	"latitude": 54.2745,
  	"longitude": 9.8958,
  	"current": ["temperature_2m",  "relative_humidity_2m"],
  	"timezone": "auto",
  	"forecast_days": 1
  }
  responses = openmeteo.weather_api(url, params=params)
  
  # Process first location. Add a for-loop for multiple locations or weather models
  response = responses[0]
  print(f"Coordinates {response.Latitude()}°N {response.Longitude()}°E")
  print(f"Elevation {response.Elevation()} m asl")
  print(f"Timezone {response.Timezone()}{response.TimezoneAbbreviation()}")
  print(f"Timezone difference to GMT+0 {response.UtcOffsetSeconds()} s")
  
  # Current values. The order of variables needs to be the same as requested.
  current = response.Current()
  current_temperature_2m = current.Variables(0).Value()
  current_relative_humidity_2m = current.Variables(1).Value()
  
  print(f"Current time {current.Time()}")
  print(f"Current temperature_2m {current_temperature_2m}")
  print(f"Current humidity_2m {current_relative_humidity_2m}")
  
  # Load .env file
  load_dotenv(dotenv_path=Path(".env"))
  
  # Get environment variables
  api_key = os.getenv("API_KEY")
  ssl_cert_path = os.getenv("SSL_CERT_PATH")
  ssl_cert_file = os.getenv("SSL_CERT_FILE")
  ssl_key_file = os.getenv("SSL_KEY_FILE")
  
  if api_key is None:
    sys.exit("API_KEY is missing")
  
  if ssl_cert_path is None or ssl_cert_file is None or ssl_key_file is None:
    sys.exit("One of the following variables is missing: SSL_CERT_PATH, SSL_CERT_FILE, SSL_KEY_FILE")
  
  
  cert = os.path.join(ssl_cert_path, ssl_cert_file)
  key = os.path.join(ssl_cert_path, ssl_key_file)
  
  # API endpoint
  url = "https://middleware:8443/climate"  # change to your actual middleware URL
  
  # Example data to send
  data = {
    "temperature": current_temperature_2m,
    "humidity": current_relative_humidity_2m,
    "location": "Westensee",
    "timestamp": datetime.now(ZoneInfo("Europe/Berlin")).isoformat()
  }
  
  # Send POST request with client certificate and headers
  response = requests.post(
      url,
      json=data,
      headers={
          "X-API-Key": f"{api_key}"
      },
      cert=(cert, key),     # client cert/key
      verify=False# optionally path to CA bundle, or False to skip verification
  )
  
  print(f"Status: {response.status_code}")
  time.sleep(10)
