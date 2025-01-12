import customtkinter as ctk
import subprocess
import psutil
from login_screen import show_login_screen

def validate_ip(ip_string):
    return True

def cleanup_process():
    """Kill all processes named 'sender.exe'"""
    for proc in psutil.process_iter(['pid', 'name']):
        try:
            # Check if process name contains 'sender.exe'
            if 'sender.exe' in proc.info['name'].lower():
                process = psutil.Process(proc.info['pid'])
                process.terminate()  # Try graceful termination first
                
                # Wait for the process to terminate
                try:
                    process.wait(timeout=3)  # Wait up to 3 seconds
                except psutil.TimeoutExpired:
                    process.kill()  # Force kill if graceful termination fails
                
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            continue

def show_ip_config_screen(app, cute_photo):

    # Create main frame
    frame = ctk.CTkFrame(master=app)
    frame.pack(pady=20, padx=20, fill="both", expand=True)

    # Title
    title = ctk.CTkLabel(frame, text="Network Configuration", font=("Roboto", 24))
    title.pack(pady=20)

    # Create entry fields for IPs
    ip_entries = {}
    
    # Node IP entries
    for i in range(1, 4):
        node_frame = ctk.CTkFrame(frame)
        node_frame.pack(pady=10, padx=20, fill="x")
        
        label = ctk.CTkLabel(node_frame, text=f"Node {i} IP:")
        label.pack(side="left", padx=10)
        
        entry = ctk.CTkEntry(node_frame, placeholder_text=f"Enter Node {i} IP")
        entry.pack(side="left", expand=True, padx=10)
        
        ip_entries[f"node{i}"] = entry

    # Server IP entry
    server_frame = ctk.CTkFrame(frame)
    server_frame.pack(pady=10, padx=20, fill="x")
    
    server_label = ctk.CTkLabel(server_frame, text="Server IP:")
    server_label.pack(side="left", padx=10)
    
    server_entry = ctk.CTkEntry(server_frame, placeholder_text="Enter Server IP")
    server_entry.pack(side="left", expand=True, padx=10)
    ip_entries["server"] = server_entry

    # Status message
    status_label = ctk.CTkLabel(frame, text="", text_color="red")
    status_label.pack(pady=10)

    def launch_sender():
        cleanup_process()

        # Validate all IPs
        ips = {
            key: entry.get().strip() 
            for key, entry in ip_entries.items()
        }
        
        # Check if any fields are empty
        if any(not ip for ip in ips.values()):
            status_label.configure(text="Please fill in all IP addresses")
            return

        # Validate IP format
        invalid_ips = [key for key, ip in ips.items() if not validate_ip(ip)]
        if invalid_ips:
            status_label.configure(
                text=f"Invalid IP format for: {', '.join(invalid_ips)}"
            )
            return

        # Construct command arguments
        args = [
            "./sender.exe",
            f"-node1={ips['node1']}",
            f"-node2={ips['node2']}",
            f"-node3={ips['node3']}",
            f"-server={ips['server']}"
        ]

        try:
            # Launch sender.exe with arguments
            process = subprocess.Popen(
                args,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            status_label.configure(
                text="Sender launched successfully",
                text_color="green"
            )
            show_login_screen(app, cute_photo)
        except Exception as e:
            status_label.configure(
                text=f"Error launching sender: {str(e)}"
            )

    # Launch button
    launch_btn = ctk.CTkButton(
        frame,
        text="Launch Sender",
        command=launch_sender
    )
    launch_btn.pack(pady=20)

    return frame