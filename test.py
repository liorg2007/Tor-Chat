from flask import Flask, request

app = Flask(__name__)

@app.route("/", methods=["GET"])
def hello_world():
    return "Hello, World!", 200

@app.route("/get-aes", methods=["POST"])
def get_aes():
    # Print request headers
    print("Headers:")
    for header, value in request.headers.items():
        print(f"  {header}: {value}")
    
    # Print request body
    body = request.data.decode('utf-8')  # Decoding the body
    print("Body:")
    print(body)

    # Respond with a success message
    return "Request received on /get-aes", 200

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=1234)
