#!/usr/bin/python3
import sys
import base64
import numpy as np
data_source = __import__(sys.argv[1])

# argv1[1] is the name of the data source
# argv[2] is input, which should be encoded
# argv[3] is expected output, which should be encoded

def data_source_res():
    return data_source.main()

def compare_result():
    data_source_result = data_source_res()
    # convert base64 to string
    expected_result = base64.b64decode(sys.argv[3]).decode("utf-8")
    try:
        deviation = data_source_result - float(expected_result)
        if deviation < 10000:
            return data_source_result
        else:
            return "null"
    except ValueError:
        return "null"

def main():
    return compare_result()

if __name__ == "__main__":
    try:
        print(main())
    except ArithmeticError:
        print("null")