#!/usr/bin/python3
import requests
import sys


def main():
    image = open(sys.argv[1], 'rb')
    r = requests.post(
        'https://image-classification.v-chain.vn/v1/short-classification',
        data={'model': 'vgg11_bn'},
        files={'image': image}
    ).json()

    if r['code'] == 200:
        return r['data'].strip()

    return None


if __name__ == "__main__":
    print(main())
