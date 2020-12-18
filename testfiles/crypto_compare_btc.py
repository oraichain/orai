#!/usr/bin/python3
import requests


def main():
    r = requests.get(
        'https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD').json()['USD']
    return r


if __name__ == "__main__":
    print(main())
