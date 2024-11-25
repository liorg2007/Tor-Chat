import customtkinter as ctk
import login_screen

def show_signup_screen(app, cute_photo):
    for widget in app.winfo_children():
        widget.destroy()

    # Header
    header = ctk.CTkFrame(master=app, fg_color="black", height=60, corner_radius=0)
    header.pack(fill="x")
    
    header_label = ctk.CTkLabel(
        master=header,
        text="M4R5HME110 S1GN-UP",
        font=("Consolas", 24, "bold"),
        text_color="green",
    )
    header_label.pack(side="top", pady=10)

    # Back Button
    back_button = ctk.CTkButton(
        master=header,
        text="< Back",
        font=("Courier New", 14),
        fg_color="black",
        text_color="green",
        hover_color="darkgreen",
        command=lambda: login_screen.show_login_screen(app, cute_photo),
    )
    back_button.pack(side="left", padx=10, pady=10)

    # Main layout frame
    main_frame = ctk.CTkFrame(master=app, fg_color="black")
    main_frame.pack(fill="both", expand=True, padx=20, pady=20)

    # Left-side signup form
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

    # Password Entry
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

    # Confirm Password Entry
    confirm_password_label = ctk.CTkLabel(
        master=form_frame,
        text="Confirm Password:",
        font=("Courier New", 16),
        text_color="green",
    )
    confirm_password_label.pack(pady=10)

    confirm_password_entry = ctk.CTkEntry(
        master=form_frame,
        width=250,
        show="*",
        fg_color="black",
        text_color="green",
        font=("Courier New", 14),
        placeholder_text="Repeat your password",
        placeholder_text_color="gray",
    )
    confirm_password_entry.pack(pady=10)

    # Data Security Checkbox
    data_security_checkbox = ctk.CTkCheckBox(
        master=form_frame,
        text="Allow data security",
        font=("Courier New", 14),
        fg_color="green",
        text_color="green",
    )
    data_security_checkbox.pack(pady=10)

    # Sign-Up Button
    signup_button = ctk.CTkButton(
        master=form_frame,
        text="Sign-Up",
        font=("Courier New", 14),
        fg_color="black",
        text_color="green",
        hover_color="darkgreen",
        command=lambda: signup_action(username_entry, password_entry, confirm_password_entry, data_security_checkbox),
    )
    # Right-side logo
    logo_label = ctk.CTkLabel(master=main_frame, image=cute_photo, text="")
    logo_label.place(relx=0.7, rely=0.5, anchor="center")
    signup_button.pack(pady=20)

# Signup action
def signup_action(username_entry, password_entry, confirm_password_entry, checkbox):
    username = username_entry.get()
    password = password_entry.get()
    confirm_password = confirm_password_entry.get()
    data_security_enabled = checkbox.get()

    if password != confirm_password:
        print("Error: Passwords do not match!")
    else:
        print(f"Sign-Up successful! Username: {username}, Data Security: {data_security_enabled}")
