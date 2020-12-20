#!/usr/bin/python3
import requests
import sys


def main():
    image = open(sys.argv[1], 'rb')

    r = requests.post(
        'http://104.248.99.206:8080/v1/identify',
        files={'image': image}
    ).json()

    print(r)
    status = r['status']

    user_id = r['data']['user_id'] if 'user_id' in r['data'] else ""

    print(status)
    status_trim = status.strip()
    status_correct = "success"
    print(status_correct)
    print(status_trim)

    # check status code of the request
    if status_trim == status_correct:
        return user_id.strip()

    return None


if __name__ == "__main__":
    print(main())
