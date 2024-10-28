#Makefile

all: main

# main: main.o queue.o
# 	g++ -o main main.o  queue.o

main: main.o
	g++ -o main main.o

main.o: main.cpp
	g++ -c main.cpp

clean:
	rm -f main *.o *~