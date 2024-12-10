FROM python:alpine3.21

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any needed dependencies specified in requirements.txt
RUN pip install --no-cache-dir "fastapi[standard]" databases asyncpg

EXPOSE 8000

# Run app.py when the container launches
CMD ["fastapi", "run", "main.py", "--port", "8000"]