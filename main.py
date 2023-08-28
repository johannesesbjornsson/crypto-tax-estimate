
import configparser
from binance.client import Client
import csv
import os


def get_binance_transactions(file):
    earn_titels = ["Simple Earn Locked Rewards", "Simple Earn Flexible Interest", "ETH 2.0 Staking Rewards"]
    tokens = {}
    with open(file) as csv_file:
        rows = csv.reader(csv_file, delimiter=',')
        for line in rows:
            trade_id, trade_date, trade_type, title, token, amount, description = line
            if title in earn_titels:
                if not token in tokens:
                    tokens[token]= {
                        trade_date.split()[0]: amount
                    }
                else:
                    tokens[token][trade_date.split()[0]] = amount
    return tokens

def main():
    client = Client(os.environ['B_API_KEY'],os.environ['B_API_SECRET'])
    config = configparser.ConfigParser()
    config.read(r'config.txt')
    tokens = get_binance_transactions(config.get('binance', 'file'))
    #for k,v in tokens["BUSD"].items():
    #    print(k,v)
    print(tokens.keys())

if __name__ == "__main__":
    main()