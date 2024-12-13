import customtkinter as ctk
import requests, json
from PIL import Image, ImageTk
from signup_screen import show_signup_screen  # Import signup screen function
from chat import open_chat_screen
from mbox import show_message_box

def show_login_screen(app, cute_photo):
    for widget in app.winfo_children():
        widget.destroy()

    # Header
    header = ctk.CTkFrame(master=app, fg_color="black", height=60, corner_radius=0)
    header.pack(fill="x")
    
    header_label = ctk.CTkLabel(
        master=header,
        text="M4R5HME110 L0G1N",
        font=("Consolas", 24, "bold"),
        text_color="green",
    )
    header_label.pack(side="top", pady=10)

    # Main layout frame
    main_frame = ctk.CTkFrame(master=app, fg_color="black")
    main_frame.pack(fill="both", expand=True, padx=20, pady=20)

    # Left-side login form
    form_frame = ctk.CTkFrame(master=main_frame, fg_color="black", width=300, height=400, corner_radius=10)
    form_frame.place(relx=0.2, rely=0.5, anchor="center")

    # Username Label and Entry
    username_label = ctk.CTkLabel(
        master=form_frame,
        text="Username:",
        font=("Courier New", 16),
        text_color="green",
    )
    username_label.pack(pady=10)

    username_entry = ctk.CTkEntry(
        master=form_frame,
        width=250,
        fg_color="black",
        text_color="green",
        font=("Courier New", 14),
        placeholder_text="Enter your username",
        placeholder_text_color="gray",
    )
    username_entry.pack(pady=10)

    # Password Label and Entry
    password_label = ctk.CTkLabel(
        master=form_frame,
        text="Password:",
        font=("Courier New", 16),
        text_color="green",
    )
    password_label.pack(pady=10)

    password_entry = ctk.CTkEntry(
        master=form_frame,
        width=250,
        show="*",
        fg_color="black",
        text_color="green",
        font=("Courier New", 14),
        placeholder_text="Enter your password",
        placeholder_text_color="gray",
    )
    password_entry.pack(pady=10)

    # Login and Sign-Up Buttons
    button_frame = ctk.CTkFrame(master=form_frame, fg_color="black")
    button_frame.pack(pady=20)

    login_button = ctk.CTkButton(
        master=button_frame,
        text="Login",
        font=("Courier New", 14),
        fg_color="black",
        text_color="green",
        hover_color="darkgreen",
        command=lambda: login_action(username_entry, password_entry, app),
    )
    login_button.pack(side="left", padx=10)

    signup_button = ctk.CTkButton(
        master=button_frame,
        text="Create User",
        font=("Courier New", 14),
        fg_color="black",
        text_color="green",
        hover_color="darkgreen",
        command=lambda: show_signup_screen(app, cute_photo),
    )
    signup_button.pack(side="left", padx=10)

    # Right-side logo
    logo_label = ctk.CTkLabel(master=main_frame, image=cute_photo, text="")
    logo_label.place(relx=0.7, rely=0.5, anchor="center")

def login_action(username_entry, password_entry, app):
    username = username_entry.get()
    password = password_entry.get()
    ans = requests.post("http://localhost:1234/login", json={"Username": username, "Password": password})
    print(f"Attempted login with username: {username} and password: {password}")
    #print(f"Received answer: {ans.}")
    try:
        if ans.status_code == 200:
            open_chat_screen(app, username)
            show_message_box(app, "Success", "You logged in!")
        else:
            json_data = json.loads(ans.content.decode())
            show_message_box(app, "Login Failed", f"{json_data['detail']['detail']}")
    except Exception as e:
        show_message_box(app, "Connection Error", f"Error: {e}")


