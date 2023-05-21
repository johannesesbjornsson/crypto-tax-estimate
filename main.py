
import configparser
from binance.client import Client
import csv
import os
import json

def get_binance_transactions(file):
    earn_titels = ["Simple Earn Locked Rewards", "Simple Earn Flexible Interest", "ETH 2.0 Staking Rewards", "Staking Rewards"]
    tokens = {}
    with open(file) as csv_file:
        rows = csv.reader(csv_file, delimiter=',')
        for line in rows:
            trade_id, trade_date, trade_type, title, token, amount, description = line
            trade_date_YYYY_MM_DD= trade_date.split()[0]
            if title in earn_titels:
                if not token in tokens:
                    tokens[token]= {
                        trade_date_YYYY_MM_DD: float(amount)
                    }
                else:
                    if trade_date_YYYY_MM_DD in tokens[token]:
                        tokens[token][trade_date_YYYY_MM_DD] += float(amount)
                    else:
                        tokens[token][trade_date_YYYY_MM_DD] = float(amount)
    return tokens

def get_gbp_value(date, token, amount):
    with open(f'historical_data/{token}.json') as f:
        data = json.load(f)
    try:
        gbp_amount = float(amount) * data[date]["GBP_PRICE"]
    except Exception as e:
        print(date, token)
        raise e
    return round(gbp_amount, 5)

def main():
    client = Client(os.environ['B_API_KEY'],os.environ['B_API_SECRET'])
    config = configparser.ConfigParser()
    config.read(r'config.txt')
    tokens = get_binance_transactions(config.get('binance', 'file'))
    columns = list(tokens.keys())


    sum_of_all_tokens = []
    for token in tokens:
        token_all_transaction = [] 
        for date, amount in tokens[token].items():
            column_index=columns.index(token)
            gbp_value = get_gbp_value(date, token, amount)
            token_all_transaction.append(gbp_value)
                
        
        sum_token = sum(token_all_transaction)
        sum_of_all_tokens.append(sum_token)
        print(token,sum_token)
        print("----------------")

    print(sum(sum_of_all_tokens))
    

if __name__ == "__main__":
    main()