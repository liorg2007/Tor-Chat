# Marshmello Space
MarshmelloSpace is a global chat that utilizes a sophisticated tor network to provide the network with top level security.


## What the project consists of

The whole project will be implemented with <span  style="color:aqua">Golang</span> and <span  style="color:green">Python</span> languages. 

There are 3 kinds of entities in the project:
- Client: Runs on Windows with a GUI
- Server: Runs on Linux system 
- Tor network node: Runs on Linux system

To simplify the Linux device deployment the project would utilize the power of Docker containers.

## Building
- The project provides docker-compose file, its enough to set up the nodes and the server
- The python GUI can be executed immediatley with the Makefile