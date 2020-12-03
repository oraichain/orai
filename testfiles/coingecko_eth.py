#!/usr/bin/python3
import requests

def main():
    r = requests.get('https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd').json()['ethereum']['usd']
    return r

if __name__ == "__main__":
    print(main())