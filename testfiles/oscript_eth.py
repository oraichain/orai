#!/usr/bin/python3
import sys
import base64
import numpy as np

# argv1[1] is the name of the data source
# argv[2] is input, which should be encoded
# argv[3] is expected output, which should be encoded

def get_data_source():
    return "crypto_compare_eth coingecko_eth"

def get_test_case():
    return "testcase_price"

def main():
    if sys.argv[1] == "aiDataSource":
        return get_data_source()
    elif sys.argv[1] == "testcase":
        return get_test_case()
    elif sys.argv[1] == "aggregation":
        results = sys.argv[2].split("-")
        try:
            results = list(map(float, results))
            aggregated_result = sum(map(float,results))
            print("aggregated result: ", aggregated_result)
            return aggregated_result / len(results)
        except ValueError:
            print("cannot convert the results into numbers")
            return 0

if __name__ == "__main__":
    try:
        print(main())
    except ArithmeticError:
        print("0")