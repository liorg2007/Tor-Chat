import os
import shutil
import PyInstaller.__main__

# Paths
project_name = "MarshmelloSpace"  # Name of the executable
entry_file = "main.py"  # Entry point of the application
output_dir = "dist"  # Directory for the built executable
icon_path = "img/icon.ico"  # Path to app icon (optional)

site_packages = os.path.join(os.path.dirname(__file__), "venv", "Lib", "site-packages")
build_output_dir = os.path.join(output_dir, project_name)  # Path to the built application folder


PyInstaller.__main__.run([
    entry_file,
    f"--name={project_name}",
    f"--distpath={output_dir}",
    "--onedir",  # Create a folder for the executable and dependencies
    "--windowed",  # Hide the console window
    f"--add-data=img{os.pathsep}img",  # Include the assets directory
    f"--icon={icon_path}",  # Add the application icon (optional)
    f"--paths={site_packages}"  # Add the virtual environment's site-packages
])

img_src = os.path.join(os.path.dirname(__file__), "img")
img_dest = os.path.join(build_output_dir, "img")

try:
    shutil.copytree(img_src, img_dest, dirs_exist_ok=True)
    print(f"Successfully copied 'img' directory to {img_dest}")
except Exception as e:
    print(f"Error copying 'img' directory: {e}")
