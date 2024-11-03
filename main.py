from transactions import TransactionHistory

def main():

    transaction_history = TransactionHistory()
    #transaction_history.load_transactions(["binance/2023_tax_year_binance.csv"])

    transaction_history.load_transactions(["cardano/Yoroi_AVA_2023-2024.csv"])
    transaction_history.load_transactions(["kraken/kraken_2023-2024.csv"])
    print(transaction_history.get_earn_transactions_value("ADA"))
    print(transaction_history.get_earn_transactions_value("MATIC"))
    print(transaction_history.get_earn_transactions_value("SOL"))
    print(transaction_history.get_earn_transactions_value("ETH"))
    print(transaction_history.get_earn_transactions_value("DOT"))



if __name__ == "__main__":
    main()

