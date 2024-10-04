# Marshmello Space
MarshmelloSpace is a global chat that utilizes a sophisticated tor network to provide the network with top level security.


## What the project consists of

The whole project will be implemented with <span  style="color:aqua">Golang</span> language. 

There are 3 kinds of entities in the project:
- Client: Runs on Windows with a GUI 
- Server: Runs on Linux system 
- Tor network node: Runs on Linux system

To simplify the Linux device deployment the project would utilize the power of Docker containers.

### Package Implementation
The project implements 3 main packagees that each entity uses:
- encryption: By definition this package provides API for encrypting in AES and RSA
- networking: Provides API for reading data and writing for TCP socket, insuring the data is passed in base64
- torpacket: Provides API for serialization and deserialization of the TOR protocol packets

Theses major packages are the core of each entity in the network, which will use it's logic to utilize the packages for the entity's responsibilities.