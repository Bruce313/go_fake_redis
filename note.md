1. listen让内核监听，epoll只是批处理，不能代listen
2. reader, writer之所以要关闭，是要close背后的文件描述符
