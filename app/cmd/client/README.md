# Building
To build the client you build two steps:
data - run in app directory the following: ``` go build -o cmd/client/data/communicator.go cmd/client/data/creator.go cmd/client/data/main.go cmd/client/data/sender.go cmd/client/data/user-requests.go ```

ui - (after installing the required libraries from requirements.txt) run python build.py
Then inside dist/my_customtkinter_app you can create shortcut from my_customtkinter_app.exe