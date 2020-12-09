import umsgpack
from flask import Flask
from flask import request
import json
app = Flask(__name__)

@app.route('/', methods=['POST'])
def hello_world():
    req_data = request.data
    print("request: ", req_data)
    unpack = json.loads(req_data)
    print("unpack: ", umsgpack.unpackb(bytearray(unpack['form']['data'])))
    return 'Hello, World!'

if __name__ == '__main__':
    app.run()