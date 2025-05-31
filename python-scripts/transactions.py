import csv
from datetime import datetime, date, timedelta
import json
import re



class TransactionHistory(object):
    def __init__(self):
        self.historical_data = {}
        for token in [ "ADA", "ARB", "AVAX", "BTC", "DOT", "ETH", "GBPUSD", "MATIC", "NEAR", "OP", "SOL"]:
            with open(f'historical_data/{token}.json') as f:
                self.historical_data[token] = json.load(f)   
        self.transactions = []

    
    def load_transactions(self, file_paths):
        for file in file_paths:
            transaction_venue = file.split("/")[0]
            with open(file) as csv_file:
                rows = csv.DictReader(csv_file, delimiter=',')
                for row in rows:
                    if transaction_venue == "binance":
                        transaction = BinanceTransaction(row, self.historical_data)
                        self.transactions.append(transaction)
                    elif transaction_venue == "cardano":
                        transaction = CardanoTransaction(row, self.historical_data)
                        self.transactions.append(transaction)
                    elif transaction_venue == "kraken":
                        transaction = KrakenTransaction(row, self.historical_data)
                        self.transactions.append(transaction)

                        
                        

    def get_earn_transactions_value(self, token):
        total_value = 0
        total_txns = 0
        for transaction in self.transactions:
            if transaction.token == token:
                total_value += transaction.value
                total_txns += 1

        return token, total_value, total_txns



class Transaction(object):
    def __init__(self, historical_data):
        self.earn_titels = None
        self.token = None
        self.transaction_date = None
        self.historical_data = historical_data
 

    def get_gransaction_type(self):
        return  self.transaction_type

    def get_transaction_id(self):
        return self.transaction_id

    def __str__(self):
        return self.transaction_id

    def get_gbp_value(self, token, amount):
        try:
            price =  self.historical_data[token][str(self.transaction_date)]["GBP_PRICE"]
            gbp_amount = float(amount) * price
        except Exception as e:
            print(self.transaction_date, token)
            raise e
        if gbp_amount < 0:
            print(self.transaction_date, token, price, float(amount) )
            raise ValueError("GBP value is negative")

        return round(gbp_amount, 5)

    def get_gbp_by_price(self, token):
        data = self.historical_data[token]
        try:
            price = data[str(self.transaction_date)]["GBP_PRICE"]
        except Exception as e:
            print(self.transaction_date, token)
            raise e
        return round(price, 5)

class BinanceTransaction(Transaction):
    def __init__(self, row, historical_data):
        super().__init__(historical_data)
        self.earn_titels = ["Simple Earn Locked Rewards", "Simple Earn Flexible Interest", "ETH 2.0 Staking Rewards", "Staking Rewards"]
        self.load_transaction(row)

    def load_transaction(self, values):
        if "Operation" in values and any(key in values["Operation"] for key in self.earn_titels):
            self.transaction_type = "earn"
            self.load_earn_transaction(values)
        elif "Side" in values:
            self.transaction_type = "trade"
        else:
            pass
            #print("Transaction type not found", values)
            #raise ValueError("Transaction type not found", values)
        pass

    def load_earn_transaction(self, values):
        self.token = values["Coin"]
        self.transaction_id = f'{values["UTC_Time"]}-{values["Coin"]}-{values["Change"]}'
        self.transaction_date = datetime.strptime(values["UTC_Time"], "%Y-%m-%d %H:%M:%S").date()

        self.value = self.get_gbp_value(values["Coin"], values["Change"])
        self.buy_price = self.get_gbp_by_price(values["Coin"])
        

    def load_trade_transaction(self, values):
        print("Not implemented yet")

class CoinbaseTransaction(Transaction):
    def __init__(self, row, historical_data):
        super().__init__(historical_data)
        super(Transaction, self).__init__()
        self.earn_titels = ["TBD"]
        self.load_transaction(row)

    def load_transaction(self, values):
        pass

class CardanoTransaction(Transaction):
    def __init__(self, row, historical_data):
        super().__init__(historical_data)
        super(Transaction, self).__init__()
        self.earn_titels = ["Staking Reward"]
        self.load_transaction(row)

    def load_transaction(self, values):
        if "Comment (optional)" in values and any(key in values["Comment (optional)"] for key in self.earn_titels):
            self.transaction_type = "earn"
            self.load_earn_transaction(values)
        else:
            pass
            #print("Transaction type not found", values)
            #raise ValueError("Transaction type not found", values)
        pass

    def load_earn_transaction(self, values):
        self.token = values["Buy Cur."]
        self.transaction_id = f'{values["Date"]}'
        self.transaction_date = datetime.strptime(values["Date"], "%Y-%m-%d %H:%M:%S").date()

        self.value = self.get_gbp_value(values["Buy Cur."], values["Buy Amount"])
        self.buy_price = self.get_gbp_by_price(values["Buy Cur."])

class KrakenTransaction(Transaction):
    def __init__(self, row, historical_data):
        super().__init__(historical_data)
        super(Transaction, self).__init__()
        self.earn_titels = ["earn","staking"]
        self.load_transaction(row)


    def load_transaction(self, values):
        #print(values)
        if values["type"] in self.earn_titels and values["subtype"] != "migration":
            self.transaction_type = "earn"
            self.load_earn_transaction(values)
        elif values["type"] == "trade":
            pass
        else:
            pass
            #print("Transaction type not found", values)
            #raise ValueError("Transaction type not found", values)
        pass

    def load_earn_transaction(self, values):
        token = re.search('([A-Z]{1,5}).*', values["asset"]).group(1)
        self.token = token
        self.transaction_id = values["txid"]
        self.transaction_date = datetime.strptime(values["time"], "%Y-%m-%d %H:%M:%S").date()
        self.value = self.get_gbp_value(token, values["amount"])
        self.buy_price = self.get_gbp_by_price(token)