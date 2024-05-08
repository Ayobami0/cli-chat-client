CLIENT  CLIENT
     \  /
    SERVER
      |
     DB

- CHATROOMS -> MULTIPLE CLIENTS
- DIRECT MESSAGE -> 2 CLIENTS
- BIDIRECTION MESSAGE
- LOGIN -> PASSWORD & USERNAME
- CREATE -> PASSWORD & USERNAME

cli-chat --help|-h
cli-chat --version|-v
cli-chat --login --username|-u=Ayobami
Enter password> ********
cli-chat --create-account
username> Ayobami
password> ********

## CLIENT TUI
+f1--------------+f2----------------------------+f3-------+
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |                              |         |
|                |f4----------------------------+---------+
|                |                          |ctrl-? - help|
+---------------------------------------------------------+

f1 -> switch to chat list
f2 -> switch to chatroom
f3 -> chatroom status
f4 -> send message
ctrl-? -> show help

## SERVER SETUP
client message ----> server ----> client response
                     |
               message database
