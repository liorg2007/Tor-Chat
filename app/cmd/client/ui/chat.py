from dateutil import parser
import threading
import time
import customtkinter as ctk
from PIL import Image, ImageTk
import requests

def format_timestamp(timestamp_str):
    """Convert ISO timestamp to a more readable format."""
    try:
        # Parse the ISO timestamp using dateutil.parser
        dt = parser.isoparse(timestamp_str)

        # Convert to local time
        local_dt = dt.astimezone()

        # Format the time
        return local_dt.strftime("%I:%M %p")
    except Exception as e:
        print(f"Error formatting timestamp: {e}")
        return timestamp_str

def start_message_receiver(chat_text: ctk.CTkTextbox, username):
    # Track the last processed message timestamp to avoid duplicates
    last_message_time = None
    
    def receive_messages():
        nonlocal last_message_time

        while True:
            try:
                # Fetch messages
                ans = requests.get("http://localhost:1234/receive-messages")
                
                if ans.status_code == 200:
                    # Parse the messages
                    messages = ans.json()

                    # Get current scroll position
                    current_scroll = chat_text.yview()[0]

                    # Clear the chat text box and insert the welcome message
                    chat_text.delete("1.0", "end")
                    chat_text.insert("0.0", f"[INFO] Welcome to the Marshmello Space, {username}! Start typing...\n")

                    # Display all messages
                    for msg in messages:
                        msg_username = msg.get('username', 'UNKNOWN')
                        msg_text = msg.get('message', '')
                        msg_time = msg.get('createdAt', '')

                        # Format timestamp
                        formatted_time = format_timestamp(msg_time)

                        if msg_username == username:
                            chat_text.insert("end", f"[YOU @ {formatted_time}] {msg_text}\n")
                        else:
                            chat_text.insert("end", f"[{msg_username} @ {formatted_time}] {msg_text}\n")

                    # Restore scroll position
                    chat_text.yview_moveto(current_scroll)

                # Wait for 1 second before the next request
                time.sleep(1)

            except Exception as e:
                print(f"Error receiving messages: {e}")
                time.sleep(5)  # Wait longer if there's an error

    # Create and start the thread
    message_thread = threading.Thread(target=receive_messages, daemon=True)
    message_thread.start()
    return message_thread

def open_chat_screen(app, username):
    # Clear existing widgets
    for widget in app.winfo_children():
        widget.destroy()

    ctk.set_appearance_mode("dark")
    ctk.set_default_color_theme("green")  # Set hacker-like green accent

    # Load images
    cute_img = Image.open("img/cute.png").resize((60, 60))
    cute_photo = ImageTk.PhotoImage(cute_img)

    # Header
    header = ctk.CTkFrame(master=app, fg_color="black", height=60, corner_radius=0)
    header.pack(fill="x")

    # Cute logo on the left
    cute_label = ctk.CTkLabel(master=header, image=cute_photo, text="")
    cute_label.pack(side="left", padx=15, pady=5)

    # Centered header title
    header_label = ctk.CTkLabel(
        master=header,
        text=f"M4R5HME110 SPACE - {username}",
        font=("Consolas", 24, "bold"),
        text_color="green",
    )
    header_label.pack(side="top", pady=10)

    # Chat section - use fill and expand to utilize whole screen
    chat_frame = ctk.CTkFrame(master=app, fg_color="black")
    chat_frame.pack(fill="both", expand=True, padx=20, pady=10)

    chat_text = ctk.CTkTextbox(
        master=chat_frame,
        wrap="word",
        fg_color="black",  # Match the hacker aesthetic
        text_color="green",  # Bright green text
        font=("Courier New", 16),  # Monospaced hacker font
        border_width=0,
    )
    chat_text.pack(fill="both", expand=True, padx=10, pady=10)
    chat_text.insert("0.0", f"[INFO] Welcome to the Marshmello Space, {username}! Start typing...\n")

    # Input field and send button
    input_frame = ctk.CTkFrame(master=app, fg_color="black")
    input_frame.pack(fill="x", pady=10, padx=20)

    input_field = ctk.CTkEntry(
        master=input_frame,
        width=640,
        placeholder_text="Type your message here...",
        font=("Courier New", 14),
        fg_color="black",
        text_color="green",
        placeholder_text_color="gray",
    )
    input_field.pack(side="left", padx=10, pady=10, expand=True, fill="x")

    send_button = ctk.CTkButton(
        master=input_frame,
        text=">>",
        width=80,
        font=("Courier New", 14),
        fg_color="black",
        text_color="green",
        command=lambda: send_message(chat_text, input_field),
    )
    send_button.pack(side="left", padx=10)

    # Ensure the cute_photo is not garbage collected
    send_button.cute_photo = cute_photo
    r = start_message_receiver(chat_text, username)

# Function to send a message
def send_message(chat_text, input_field):
    message = input_field.get()
    
    ans = requests.post("http://localhost:1234/send-message", json={"Message": message})
    
    print(f"Attempted to send a message: {message}")
    print(f"Received answer: {ans}")

    input_field.delete(0, 'end')

