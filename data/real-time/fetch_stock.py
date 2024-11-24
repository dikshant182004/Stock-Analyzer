import sys
import json
from bsedata.bse import BSE

def fetch_stock_data(symbol):
    bse = BSE(update_codes=True)
    try:
        stock_info = bse.getQuote(symbol)
        return {
            "symbol": symbol,
            "price": float(stock_info["currentValue"]),  # Convert to float
            "high": float(stock_info["dayHigh"]),
            "low": float(stock_info["dayLow"]),
            "lastUpdated": stock_info["updatedOn"],
            "error": ""
        }
    except Exception as e:
        return {"symbol": symbol, "error": str(e)}

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(json.dumps({"error": "Stock symbol not provided"}))
        sys.exit(1)

    symbol = sys.argv[1]
    result = fetch_stock_data(symbol)
    print(json.dumps(result,indent=4))
