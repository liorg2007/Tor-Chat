# Building
First create a venv inside the ui directory and switch to it.

Then install libraries using requirements.txt.

ui - (after installing the required libraries from requirements.txt) run python build.py

To build the client you build two steps:
data - run in app directory the following: ``` go build -o cmd/client/ui/dist/MarshmelloSpace/sender.exe cmd/client/data/communication.go cmd/client/data/creator.go cmd/client/data/main.go cmd/client/data/sender.go cmd/client/data/user-requests.go ```

Then inside dist/my_customtkinter_app you can create shortcut from my_customtkinter_app.exe