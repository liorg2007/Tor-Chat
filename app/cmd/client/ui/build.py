import os
import PyInstaller.__main__

# Paths
project_name = "my_customtkinter_app"  # Name of the executable
entry_file = "main.py"  # Entry point of the application
output_dir = "dist"  # Directory for the built executable
icon_path = "img/icon.ico"  # Path to app icon (optional)
customtkinter_path = os.path.join(os.path.dirname(__file__), "venv", "Lib", "site-packages", "customtkinter")  # Adjust if needed

# PyInstaller build command
PyInstaller.__main__.run([
    entry_file,
    f"--name={project_name}",
    f"--distpath={output_dir}",
    "--onedir",  # Create a folder for the executable and dependencies
    "--windowed",  # Hide the console window
    f"--add-data=img{os.pathsep}img",  # Include the assets directory
    f"--add-data={customtkinter_path}{os.pathsep}customtkinter",  # Include the CustomTkinter directory
    f"--icon={icon_path}",  # Add the application icon (optional)
])
