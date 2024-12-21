import streamlit as st
import redis
import threading
import json
import time
from collections import deque

# A deque to hold the latest data received from Redis
data_buffer = deque(maxlen=100)

# Initialize Redis client
redis_client = redis.StrictRedis(host='localhost', port=6379, db=0)

# Function to subscribe to Redis channels and listen for messages
def subscribe_to_redis():
    pubsub = redis_client.pubsub()
    pubsub.psubscribe('*')  # Subscribe to all channels of different symbols name 
    for message in pubsub.listen():
        if message['type'] == 'pmessage':
            # Append the message to the data buffer
            channel = message['channel'].decode('utf-8')
            payload = message['data'].decode('utf-8')

            # Optionally, parse JSON payload if it's structured
            try:
                payload = json.loads(payload)
            except json.JSONDecodeError:
                pass  # If not JSON, keep it as a string

            data_buffer.append((channel, payload))

# Function to display data in the Streamlit app
def display_data():
    st.title("Real-time Stock Data")
    
    # Create a placeholder for data
    placeholder = st.empty()

    while True:
        # Wait for new data to appear in the buffer
        if data_buffer:
            channel, data = data_buffer.popleft()
            st.write(f"Channel: {channel}")
            st.write(data)
        time.sleep(1)

# Start Redis subscription in a separate thread
subscription_thread = threading.Thread(target=subscribe_to_redis, daemon=True)
subscription_thread.start()

# Run the Streamlit app to display the data
display_data()
