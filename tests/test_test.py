import requests
import os
from dotenv import load_dotenv
import time
import matplotlib.pyplot as plt
import seaborn as sns
from concurrent.futures import ThreadPoolExecutor, as_completed

# Load environment variables from the .env file
load_dotenv()

# API endpoint and headers
url = 'https://127.0.0.1:8443/data'
API_KEY: str = os.getenv('API_KEY')
headers = {
    'X-API-Key': API_KEY
}

# Disable SSL verification
verify_ssl = False

# Collecting metrics
response_times = []
status_codes = []

# Define number of requests and concurrent threads
num_requests = 100  # Modify the number of requests you want to send
num_threads = 2# Modify the number of concurrent threads

def send_request():
    start_time = time.time()
    try:
        response = requests.post(url, headers=headers, verify=verify_ssl)
        response_time = time.time() - start_time
        response_times.append(response_time)
        status_codes.append(response.status_code)
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")
        response_times.append(None)
        status_codes.append(None)

# Using ThreadPoolExecutor to send requests concurrently
with ThreadPoolExecutor(max_workers=num_threads) as executor:
    futures = [executor.submit(send_request) for _ in range(num_requests)]
    # Wait for all futures to complete
    for future in as_completed(futures):
        future.result()  # This will raise exceptions if any requests failed

# Data analysis: Basic statistics
avg_response_time = sum(response_times) / len(response_times) if response_times else None
successful_requests = sum(1 for code in status_codes if code == 200)
failed_requests = sum(1 for code in status_codes if code != 200)

print(f"Average Response Time: {avg_response_time:.4f} seconds")
print(f"Successful Requests: {successful_requests}")
print(f"Failed Requests: {failed_requests}")

# Visualizing response times and status codes
plt.figure(figsize=(10, 6))

# Response time visualization (Boxplot)
plt.subplot(2, 2, 1)  # 2 rows, 2 columns, 1st plot
sns.boxplot(data=response_times)
plt.title('Response Time Distribution')
plt.ylabel('Response Time (seconds)')

# Status code visualization (Bar plot)
status_code_counts = {code: status_codes.count(code) for code in set(status_codes) if code is not None}

plt.subplot(2, 2, 2)  # 2 rows, 2 columns, 1st plot
sns.barplot(x=list(status_code_counts.keys()), y=list(status_code_counts.values()))
plt.title('API Response Status Code Distribution')
plt.xlabel('Status Code')
plt.ylabel('Frequency')

# Response Time Histogram
plt.subplot(2, 2, 3)  
sns.histplot(response_times, kde=True)
plt.title('Response Time Histogram')
plt.xlabel('Response Time (seconds)')
plt.ylabel('Frequency')
plt.tight_layout()

# Show the plot
plt.show()
