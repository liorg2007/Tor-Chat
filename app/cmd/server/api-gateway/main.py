from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import JSONResponse
import httpx

app = FastAPI()

AUTH_SERVICE = "http://auth-service:8000"
MESSAGE_SERVICE = "http://message-service:8000"

services = ['auth', 'messages']
auth_paths = ['register', 'login', 'users']
message_path = ['send', 'fetch']

'''
API gateway:
Register an account: /auth/register
Login into account: /auth/login

Send message: /messages/send
Get messages: /messages/fetch
'''

async def send_to_service(service_path: str, request: Request):
    try:
        json_data = await request.json()
    except:
        json_data = None  # Handle cases where no JSON body is sent

    async with httpx.AsyncClient() as client:
        if request.method == "POST":
            response = await client.post(service_path, json=json_data)
        elif request.method == "GET":
            response = await client.get(service_path)
        elif request.method == "DELETE":
            response = await client.delete(service_path)
        else:
            raise HTTPException(status_code=405, detail="Method not allowed")

    if response.status_code != 200:
        raise HTTPException(
            status_code=response.status_code,
            detail=response.json() if response.headers.get("content-type") == "application/json" else response.text,
        )

    try:
        return response.json()  # Return JSON data
    except ValueError:
        return response.text  # If not JSON, return plain text
    

@app.api_route("/{service}/{path:path}", methods=["GET", "POST", "DELETE"])
async def catch_all(request: Request, service:str, path: str):
    if service not in services:
        raise HTTPException(status_code=404, detail="Service not found")
    
    # Check if token validation is neccesary
    if service != "auth":
        response = ""
        token = ""
        try:
            json_data = await request.json()
            token = json_data["Token"]
        except:
            raise HTTPException(status_code=400, detail="You got to send a json")
        
        if not token:
            raise HTTPException(status_code=401, detail="Token is required for authentication")

        async with httpx.AsyncClient() as client:
            response = await client.post(f"{AUTH_SERVICE}/auth/jwt_val", json={"token": token})

        if response.status_code != 200:
            raise HTTPException(status_code=response.status_code, detail=response.json().get("detail"))
        else:
            res = response.json()
    # Now Call the right handler
    if service == 'auth':
        return await handle_auth(request, path)
    elif service == 'messages':
        return await handle_messages(request, path)

    raise HTTPException(status_code=405, detail="Method not allowed")



async def handle_messages(request: Request, path: str):
    if path not in message_path:
        raise HTTPException(status_code=404, detail="Service doesn't exist")

    try:
        return await send_to_service(MESSAGE_SERVICE + "/messages/" + path, request)
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

async def handle_auth(request: Request, path: str):
    if path not in auth_paths:
        raise HTTPException(status_code=404, detail="Service doesn't exist")

    try:
        return await send_to_service(AUTH_SERVICE + "/auth/" + path, request)
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
