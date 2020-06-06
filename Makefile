TARGET := bingo

$(TARGET): $(TARGET).go
	go build $<

.PHONY: clean
clean:
	rm -f $(TARGET)
