import streamlit as st
import json
import streamlit.components.v1 as com
from streamlit_lottie import st_lottie
from streamlit_option_menu import option_menu

st.set_page_config(
    page_title="STOCK ANALYZER",
    page_icon=":chart_with_upwards_trend:",
    layout="centered",
    initial_sidebar_state="expanded",
    menu_items={
        'Get Help': 'https://www.extremelycoolapp.com/help',
        'Report a bug': "https://www.extremelycoolapp.com/bug",
        'About': "# This is a header. This is an *extremely* cool app!"
    }
)
st.markdown("""
<style>
.stDeployButton
{
    visibility:hidden;
}       
</style>
""",unsafe_allow_html=True)

with st.sidebar:
    with open("animation.json") as source:
        animation=json.load(source)
        temp=st_lottie(animation,width=200,height=158,quality="medium",)
st.sidebar.success("**Lets Take a Tour !!!**")
with st.sidebar:
    selected=option_menu(
        menu_title="Menu",
        options=["Home","Web_scraping"],
        default_index=0,
        menu_icon="cast"
    )

if selected == "Home":
    st.markdown('<h3 style="text-align:center; color:#FF5733; font-size: 30px; font-weight: bold;"> STOCK ANALYZER </h3>',
                    unsafe_allow_html=True)
    st.markdown("<h1 style='text-align:center; font-size:25px;'>Welcome to Stock Analyzer, your daily Stock News at one place!</h1>", unsafe_allow_html=True)
    st.markdown("---")
    # st.write("<p style='font-size: 18px;'>Data Glimpse is a Streamlit application designed to automate the web scraping process. Whether you need data from e-commerce giants like Amazon, eBay, or Walmart, or wish to extract information from news articles, blogs, or forums, Data Glimpse provides a seamless solution.</p>", unsafe_allow_html=True)

    # st.write("<p style='font-size: 18px;'>With Data Glimpse, you can effortlessly gather valuable insights and information from diverse sources on the web. By simply inputting the URL of the website you want to scrape data from, and specifying the relevant tag details such as class names, IDs, or XPath, Data Glimpse automates the data extraction process.</p>", unsafe_allow_html=True)

    # st.write("<p style='font-size: 18px;'>Experience the power of web scraping with Data Glimpse now! Unleash the potential of data-driven decision-making, market analysis, and trend spotting. Gain a competitive edge by accessing real-time data from various online platforms, and utilize it for research, analytics, or business intelligence purposes.</p>", unsafe_allow_html=True)


elif selected =="Web_scraping":
    pass
    # scraper=p1.WebScraper()
    