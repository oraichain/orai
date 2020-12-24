#!/usr/bin/python3
import sys
import base64

# argv1[1] is the name of the data source
# argv[2] is input, which should be encoded
# argv[3] is expected output, which should be encoded

def main():
    if sys.argv[1] == "aggregation":
        try:
            return sys.argv[2]
        except Exception:
            return 0

if __name__ == "__main__":
    try:
        print(main())
    except ArithmeticError:
        print("0")