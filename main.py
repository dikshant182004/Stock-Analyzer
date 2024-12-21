import streamlit as st
import threading
import time

# Simulated data update function
def fetch_data():
    while True:
        time.sleep(1)
        if "data" not in st.session_state:
            st.session_state["data"] = []
        st.session_state["data"].append({"timestamp": time.time()})

# Start background thread
if "thread" not in st.session_state:
    st.session_state["thread"] = threading.Thread(target=fetch_data, daemon=True)
    st.session_state["thread"].start()

# Streamlit app UI
st.title("Real-time Stock Updates")
if "data" in st.session_state:
    for item in st.session_state["data"][-10:]:
        st.write(item)


# celery 
# Go
# streamlit
# Redis
# Postgresql
