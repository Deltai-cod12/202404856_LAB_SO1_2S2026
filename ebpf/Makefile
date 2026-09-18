CLANG ?= clang

BPF_CFLAGS := -O2 -g -target bpf \
	-D__TARGET_ARCH_x86 \
	-I/usr/include/$(shell uname -m)-linux-gnu

TARGET := kill_monitor.bpf.o
SOURCE := kill_monitor.bpf.c

all: $(TARGET)

$(TARGET): $(SOURCE) vmlinux.h
	$(CLANG) $(BPF_CFLAGS) \
		-c $(SOURCE) \
		-o $(TARGET)

clean:
	rm -f $(TARGET)