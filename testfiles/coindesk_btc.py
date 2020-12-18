#!/usr/bin/python3
import requests
# import numpy


def main():
    r = requests.get(
        'https://api.coindesk.com/v1/bpi/currentprice.json').json()['bpi']['USD']['rate_float']
    return r


if __name__ == "__main__":
    print(main())
