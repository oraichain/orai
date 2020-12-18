#!/usr/bin/python3
import requests
import sys


def main():
    image = open(sys.argv[1], 'rb')
    r = requests.post(
        'https://ocr.v-chain.vn/v1/ocr',
        files={'image': image}
    ).json()

    if r['code'] == 200:
        return r['data'].strip()

    return None


if __name__ == "__main__":
    print(main())
