import customtkinter as ctk
from PIL import Image, ImageTk

# Set up the appearance and theme
ctk.set_appearance_mode("dark")
ctk.set_default_color_theme("green")  # Set hacker-like green accent

# Create the main app window
app = ctk.CTk()
app.geometry("800x600")
app.title("Marshmello Hacker Space")

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
    text="M4R5HME110 SPACE",
    font=("Consolas", 24, "bold"),
    text_color="green",
)
header_label.pack(side="top", pady=10)

# Chat section
chat_frame = ctk.CTkScrollableFrame(master=app, fg_color="black", border_width=1, border_color="green")
chat_frame.pack(pady=20, padx=20, fill="both", expand=True)

chat_text = ctk.CTkTextbox(
    master=chat_frame,
    wrap="word",
    fg_color="black",  # Match the hacker aesthetic
    text_color="green",  # Bright green text
    font=("Courier New", 16),  # Monospaced hacker font
    border_width=0,
)
chat_text.pack(fill="both", expand=True, padx=10, pady=10)
chat_text.insert("0.0", "[INFO] Welcome to the Marshmello Space! Start typing...\n")

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
input_field.pack(side="left", padx=10, pady=10)

send_button = ctk.CTkButton(
    master=input_frame,
    text=">>",
    width=80,
    font=("Courier New", 14),
    fg_color="black",
    text_color="green",
    command=lambda: send_message(),
)
send_button.pack(side="left", padx=10)

# Function to send a message
def send_message():
    message = input_field.get()
    if message.strip():
        chat_text.insert("end", f"[USER] {message}\n")
        chat_text.see("end")
        input_field.delete(0, "end")

# Start the application
app.mainloop()
