from binance.client import Client
import os 
from datetime import datetime, date, timedelta
import json
from binance.exceptions import BinanceAPIException
import csv
import re

# Get data from https://uk.investing.com/
def get_usd_gbp_data():
    date_prices = {}
    with open('usdt_data.csv', newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            row_date = datetime.strptime(row['\ufeff"Date"'], "%d/%m/%Y").date()
            
            date_prices[str(row_date)] = {
                "PRICE": round(1/float(row["Price"]),5)
            }

    
    # Same dates as CSV downloaded
    start_date = date(2021, 4, 2)
    end_date = date(2024, 11, 1)
    # Adding missing days with last known price
    while start_date <= end_date:
        if not str(start_date) in date_prices:
            date_prices[str(start_date)] = last_price
        else: 
            last_price = date_prices[str(start_date)]
        start_date += timedelta(days=1)

    with open(f'historical_data/GBPUSD.json', 'w') as f:
        json.dump(date_prices, f, indent=4, sort_keys=True)

def calculate_synthetic_gbp_price(date, usdt_price):
    with open('historical_data/GBPUSD.json') as f:
        data = json.load(f)

    gcp_price = usdt_price / data[date]["PRICE"]
    return round(gcp_price, 5 )

def write_pair_to_file(client, pair, calculate_synthetic_price=True, filename=None,start_date="6 April, 2021", end_date="5 April, 2024"):
    klines = client.get_historical_klines(pair, Client.KLINE_INTERVAL_1DAY, start_date, end_date)
    date_prices = {}
    for k in klines:
        date = datetime.fromtimestamp(k[0]/1000).strftime('%Y-%m-%d')
        price = float(k[1])
        if calculate_synthetic_price:
            date_prices[date] = { 
                "PRICE": price,
                "GBP_PRICE": calculate_synthetic_gbp_price(date, price)
            }
        else:
            date_prices[date] = { 
                "PRICE": price
            }
    if not filename:
        filename = pair
    with open(f'historical_data/{filename}.json', 'w') as f:
        json.dump(date_prices, f, indent=4)


def main():
    client = Client(os.environ['B_API_KEY'],os.environ['B_API_SECRET'])
    #write_pair_to_file(client, "GBPUSDT", calculate_synthetic_price=False)
    #write_pair_to_file(client, "BTCGBP", calculate_synthetic_price=False)
    #write_pair_to_file(client, "ETHGBP", calculate_synthetic_price=False)

    tokens = ['BTC', 'ETH', 'ADA', 'SOL', 'DOT', 'MATIC', 'AVAX', 'NEAR', 'OP', "ARB"]
    for token in tokens:
        try:
            pair = token+"USDT"
            write_pair_to_file(client, pair, calculate_synthetic_price=True, filename=token)
        except BinanceAPIException as e:
            print(e)

if __name__ == "__main__":
    main()
    #get_usd_gbp_data()

