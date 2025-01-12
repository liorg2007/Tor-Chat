import customtkinter as ctk
from PIL import Image, ImageTk
from ip_config_screen import show_ip_config_screen
import subprocess
import threading

def launch_communicator():
    try:
        subprocess.Popen(["./sender.exe"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    except FileNotFoundError:
        print("Error: 'sender.exe' not found. Please ensure it's in the correct directory.")
    except Exception as e:
        print(f"An error occurred while launching 'sender.exe': {e}")

# Launch the communicator in a separate thread
threading.Thread(target=launch_communicator, daemon=True).start()

# Set up the appearance and theme
ctk.set_appearance_mode("dark")
ctk.set_default_color_theme("green")

# Create the main app window
app = ctk.CTk()
app.geometry("800x600")
app.title("Marshmello Hacker App")

# Load logo image
cute_img = Image.open("img/cute.png").resize((400, 400))
cute_photo = ImageTk.PhotoImage(cute_img)

# Start with login screen
show_ip_config_screen(app, cute_photo)

# Start the application
app.mainloop()
