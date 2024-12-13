import customtkinter as ctk

class CustomMessageBox(ctk.CTkToplevel):
    def __init__(self, parent, title="Message", message=""):
        super().__init__(parent)
        
        # Set properties
        self.title(title)
        self.geometry("300x150")
        self.resizable(False, False)
        self.grab_set()  # Make the message box modal
        
        # Configure the window to be centered
        self.after(10, self._center_window)
        
        # Message Label
        self.message_label = ctk.CTkLabel(self, text=message, wraplength=280, font=("Courier New", 14))
        self.message_label.pack(pady=20, padx=10, expand=True)
        
        # OK Button
        self.ok_button = ctk.CTkButton(self, text="OK", command=self.close)
        self.ok_button.pack(pady=10)
    
    def _center_window(self):
        # Center the window on the parent
        self.update_idletasks()
        parent_width = self.master.winfo_width()
        parent_height = self.master.winfo_height()
        parent_x = self.master.winfo_rootx()
        parent_y = self.master.winfo_rooty()
        
        window_width = self.winfo_width()
        window_height = self.winfo_height()
        
        x = parent_x + (parent_width - window_width) // 2
        y = parent_y + (parent_height - window_height) // 2
        
        self.geometry(f'+{x}+{y}')
    
    def close(self):
        self.grab_release()  # Release the grab
        self.destroy()

def show_message_box(parent, title="Message", message=""):
    CustomMessageBox(parent, title, message)