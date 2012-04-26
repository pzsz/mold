SRC = init.cpp
OBJ = $(SRC:.cpp=.o)

OUT = libbb.a
OUT_SO = libbb.so

# include directories
INCLUDES = -I. -I/usr/local/include/bullet

# C++ compiler flags (-g -O2 -Wall)
CCFLAGS = -g -Wall -fPIC -m64

# compiler
CCC = g++

# library paths
LIBS = -L/usr/local/lib -lm -lBulletSoftBody -lBulletDynamics -lBulletCollision -lLinearMath

# compile flags
LDFLAGS = -g

default: $(OUT_SO)

.cpp.o:
	$(CCC) $(INCLUDES) $(CCFLAGS) -c $< -o $@

$(OUT): $(OBJ)
	ar rcs $(OUT) $(OBJ)

$(OUT_SO): $(OBJ)
	$(CCC) $(CFLAGS) $(LIBS) -fPIC -shared -o $@ $^

clean:
	rm -f $(OBJ) $(OUT) $(OUT_SO) Makefile.bak 
