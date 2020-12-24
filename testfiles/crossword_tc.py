#!/usr/bin/python3
import sys
import base64
import json
# from importlib.util import spec_from_loader, module_from_spec
# from importlib.machinery import SourceFileLoader 

# # import module without using the .py extension
# spec = spec_from_loader(sys.argv[1], SourceFileLoader(sys.argv[1], "./" + sys.argv[1]))
# data_source = module_from_spec(spec)
# spec.loader.exec_module(data_source)
data_source = __import__(sys.argv[1])

# argv1[1] is the name of the data source
# argv[2] is input, which should be encoded
# argv[3] is expected output, which should be encoded

def data_source_res():
    return data_source.main()

def main():
    return data_source_res()

if __name__ == "__main__":
    try:
        print(main())
    except ArithmeticError:
        print("null")