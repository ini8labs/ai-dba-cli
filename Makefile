

# Build the project
build:
	go build -o dba.exe ./main.go

# Run the binary
run: 
	./dba.exe
	

# Clean up the binary
clean:
	rm  ./dba.exe

# Help command
help:
	@echo "Usage:"
	@echo "  make build      Build the binary"
	@echo "  make run        Build and run the binary"
	@echo "  make clean      Remove the binary"
