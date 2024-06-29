Felipe Fernandez Aguilar Rol: 201856602-0
Christian Renato Ossa Salgado Rol: 201873565-5
Rol: 201873565-5

datos maquina virtual:

    dist037 - YnazhbFcu56A
    
Instrucciones para ejecucion de la tarea:

primero se debe entrar a la maquina virtual para el desarrollo de la tarea con las credenciales que dejaremos
una vez hecho esto se puede ver la carpeta al utilizar el comando ls dentro de la maquina virtual y se debera proceder a ingresar a esta 
una vez hecho esto se deben correr los siguientes comandos en orden:

    - comando para crear la red compartida de dockers
    docker network create lab5

    - comando para construir las imagenes
    docker build -f dockerfile.broker -t broker .
    docker build -f dockerfile.comandante -t comandante .
    docker build -f dockerfile.jeth -t jeth .
    docker build -f dockerfile.malkor -t malkor .
    docker build -f dockerfile.fulcrum --build-arg NUM=1 -t fulcrum1 .
    docker build -f dockerfile.fulcrum --build-arg NUM=2 -t fulcrum2 .
    docker build -f dockerfile.fulcrum --build-arg NUM=3 -t fulcrum3 .

    - comando para correr los servidores fulcrum y el broker

    docker run -p 50056:50056 -d --network lab5 --name fulcrum1 fulcrum1
    docker run -p 50057:50057 -d --network lab5 --name fulcrum2 fulcrum2
    docker run -p 50058:50058 -d --network lab5 --name fulcrum3 fulcrum3
    docker run -p 50051:50051 -d --network lab5 --name broker broker

    - comandos para las distintas terminales abiertas en la maquina (cada una de estas debe ser ejecutada en una terminales distintas en la misma maquina)

    docker run -i -p 50050:50050 --network lab5 --name comandante comandante
    docker run -i -p 50052:50052 --network lab5 --name jeth jeth
    docker run -i -p 50053:50053 --network lab5 --name malkor malkor

Para el correcto uso de los programas se debe seguir la estructura planteada en el contexto de la tarea, sino esta fallara.

en caso de que el comandante falle con este error:
2024/06/29 03:12:08 Failed to execute command: rpc error: code = Unknown desc = stale read
hay que borrar el container y volver a correr con su comando