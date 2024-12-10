from fastapi import FastAPI, Request, HTTPException
import jwt
from datetime import datetime
from pymongo import MongoClient, ASCENDING 
app = FastAPI()

MONGO_URL = "mongodb://localhost:27017"

collection = None

async def store_message(username, message):
    global collection
    try:
        document = {
            "username": username,
            "message": message,
            "createdAt": datetime.utcnow()  # Field used by TTL index
        }
        collection.insert_one(document)
        print(f"Message for {username} stored successfully!")
    except Exception as e:
        print(str(e))
        raise HTTPException(status_code=400, detail="Can't send message")

@app.post("/message/fetch")
async def fetch_messages(request: Request):
    try:
        messages = list(collection.find({}, {"_id": 0}))
        return messages
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error fetching messages: {str(e)}")

@app.post("/message/send")
async def send_message(request: Request):
    # Extract data from JSON
    json_data = await request.json()
    message = token = username = ""
    try:
        message = json_data['message']
        token = json_data['token']
    except:
        raise HTTPException(status_code=400, detail='json fields: {"message":message, "token":token}')
    
    decoded = jwt.decode(token, options={"verify_signature": False})

    try:
        username = decoded['username']
    except:
        raise HTTPException(status_code=400, detail='Invalid token')
    
    await store_message(username, message)

    return {"status" : "success"}
    
@app.on_event("startup")
async def startup():
    global collection
    # Connect to the MongoDB container
    client = MongoClient("mongodb://message-db:27017")  # Use "mongodb://<container_name>:27017" if in Docker network
    db = client.message_db  # Database
    collection = db.messages  # Collection

    # Ensure a TTL index is created on the 'createdAt' field
    collection.create_index("createdAt", expireAfterSeconds=30)  # 30 seconds expiration
