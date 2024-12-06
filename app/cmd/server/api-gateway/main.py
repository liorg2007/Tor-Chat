from fastapi import FastAPI, Request

app = FastAPI()

AUTH_SERVICE = "http://auth-service:8000/"
MESSAGE_SERVICE = "http://message-service:8000/"

'''
API gateway:
Register an account: /auth/register
Login into account: /auth/login

Send message: /messages/send
Get messages: /messages
'''



