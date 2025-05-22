FROM python:3.12-slim

# Copy entire app folder (including requirements.txt, .env, certs, and script)
WORKDIR /app
COPY getweather/ .
COPY .env .
COPY  middleware/cert/ cert/

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
# Run the script
CMD ["python", "main.py"]

