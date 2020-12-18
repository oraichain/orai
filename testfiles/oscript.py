#!/usr/bin/python3
import requests
import sys


def aggregate(num1, num2):
    return float(num1) + float(num2)


if __name__ == "__main__":
    print(aggregate(sys.argv[1], sys.argv[2]))
