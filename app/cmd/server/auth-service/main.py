from fastapi import FastAPI, Request, HTTPException
from datetime import timedelta, datetime, timezone
from databases import Database
import os, hashlib, secrets, dotenv, jwt

app = FastAPI()

# Get the DATABASE_URL from environment variables
DATABASE_URL = os.getenv("DATABASE_URL")

# Create a Database instance
database = Database(DATABASE_URL)

# Global jwt credentials
secret = ""
algorithm = "HS256"   

# Creates a JWT session token
# Input: Username
# Output: A JWT session token for 1800
async def create_jwt(username: str) -> str:
    payload = {
    'username': username,
    'exp': datetime.now(timezone.utc) + timedelta(seconds=1800)
    }

    token = jwt.encode(payload, secret, algorithm)

    return token

# Validates a JWT token
# Input: JWT token
# Output: Returns success if the token is valid
@app.post("/auth/jwt_val")
async def validate_jwt(request: Request):
    # Extract the token from JSON
    json_data = await request.json()
    token = json_data['token']
    data = ''
    print(request)
    try: 
        # This also gets data from the token and validates it (exception if invalid)
        data = jwt.decode(token, secret, algorithms=[algorithm])
    except Exception as e:
        raise HTTPException(status_code=401, detail="Invalid token")
    
     # Check if the user even exists
    check_query = "SELECT id FROM users WHERE name = :username"
    existing_user = await database.fetch_one(query=check_query, values={"username": data["username"]})

    if not existing_user:
        # If the user exists, raise an error
        raise HTTPException(status_code=400, detail="User does'nt exist.")
    
    return {"status": "valid"}
    

# Register a new user in the system
# Input: Username and password
# Output: Return success if user doesnt exist
@app.post("/auth/register")
async def register_account(request: Request):
    # Extract JSON data
    json_data = await request.json()
    username = password = ""
    try:
        username = json_data['username']
        password = json_data['password']
    except:
        raise HTTPException(status_code=400, detail='json fields: {"usernmae":username, "passowrd":password}')

    # Hash the password (e.g., using SHA-256)
    hashed_password = hashlib.sha256(password.encode()).hexdigest()

    # Check if the user already exists
    check_query = "SELECT id FROM users WHERE name = :username"
    existing_user = await database.fetch_one(query=check_query, values={"username": username})

    if existing_user:
        # If the user exists, raise an error
        raise HTTPException(status_code=400, detail="User already exists.")

    # Insert the new user
    insert_query = """
    INSERT INTO users (name, password_hash)
    VALUES (:username, :hashed_password)
    """
    values = {"username": username, "hashed_password": hashed_password}
    await database.fetch_one(query=insert_query, values=values)

    # Return success
    return {"status": "success"}


# Authenticate as suer
# Input: Username and Password
# Output: If valid return a JWT token
@app.post("/auth/login")
async def login_account(request: Request):
    # Extract data from JSON
    json_data = await request.json()
    username = password = ""
    try:
        username = json_data['username']
        password = json_data['password']
    except:
        raise HTTPException(status_code=400, detail='json fields: {"usernmae":username, "passowrd":password}')


    # Hash the password (e.g., using SHA-256)
    password_hash = hashlib.sha256(password.encode()).hexdigest()

    # Check if the user and password exists
    check_query = "SELECT id FROM users WHERE name = :username AND password_hash = :password_hash"
    existing_user = await database.fetch_one(query=check_query, values={"username": username, "password_hash": password_hash})

    if not existing_user:
        # If the user exists, raise an error
        raise HTTPException(status_code=400, detail="Bad login credentials :(")

    # Create a session token
    jwt_token = await create_jwt(username)

    # Return the token
    return {"status": "success", "token": jwt_token} 


# Show all users in database
@app.get("/auth/users")
async def get_users(request: Request):
    insert_query = """
    SELECT * FROM users;
    """

    users = await database.fetch_all(query=insert_query)

    # Return the newly created user details
    return {"status": "success", "users": users}  # Return the list users


# Delete all users from the database
@app.delete("/auth/users")
async def clear_users(request: Request):
    insert_query = """
    DELETE FROM users;
    """

    users = await database.fetch_one(query=insert_query)

    return {"status": "success"}  

@app.on_event("startup")
async def startup():
    global secret, algorithm
    # Query to create the users table (if it doesn't already exist)
    db_init_query = """
    CREATE TABLE IF NOT EXISTS users (
        name VARCHAR(100) NOT NULL PRIMARY KEY,
        password_hash VARCHAR(255) NOT NULL
    );
    """
    await database.connect()  # Ensure the database connection is established
    await database.execute(db_init_query)  # Execute the table creation query

    # If the JWT data isn't initialzie in .env file then create it
    if(not os.path.isfile(".env")):
        secret = secrets.token_hex(20) # Create a secret value for JWT auth
        with open(".env", "w") as jwt_data:
            jwt_data.writelines(["secret = " + secret, " algorithm = HS256"])

    dotenv.load_dotenv() # Load the secret value
    secret = os.getenv('secret')


@app.on_event("shutdown")
async def shutdown():
    await database.disconnect()  # Cleanly close the database connection

    
