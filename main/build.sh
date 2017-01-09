g++ -c -fpic Demo.cpp -I./ 

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./

g++ -g -o Demo.out Demo.o -L./ -lTDBAPI 
