import customtkinter as ctk
from PIL import Image, ImageTk
from login_screen import show_login_screen

# Set up the appearance and theme
ctk.set_appearance_mode("dark")
ctk.set_default_color_theme("green")

# Create the main app window
app = ctk.CTk()
app.geometry("800x600")
app.title("Marshmello Hacker App")

# Load logo image
cute_img = Image.open("img/cute.png").resize((400, 400), Image.ANTIALIAS)
cute_photo = ImageTk.PhotoImage(cute_img)

# Start with login screen
show_login_screen(app, cute_photo)

# Start the application
app.mainloop()
