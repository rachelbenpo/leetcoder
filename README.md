# leetcoder

### this application contains server and client.

### run client:
here is the exe file:
https://drive.google.com/file/d/17LmwT1cvD88dwsfn1doQvlf76nJ-WSJp/view?usp=drive_link

#### build client from source:
"go build" in the client folder, and then run leetcode-client.exe file

you can add it to your path env. and then type "leetcode-client" in the terminal if you want


### run server:
"go build" in the server folder, and then run leetcode-server.exe file

you can add it to your path env. and then type "leetcode-server" in the terminal if you want

### credentials:
once you run the app for first time, you would be asked for your github user name, and token.
after running first time the server, go to ghcr.io/yourName/checking-container, and set it to public.
restart the server.


#### instructions of use:
the format of answer func is specified in the question instructions.

the ans function name should be the name of the question

(notice when you insert a question, to name it in correct name)

in your answer, please use "\\n" for specify a new line


you should have docker v 1.44 minimum, with k8s enabled on your computer.