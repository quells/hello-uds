CFLAGS=-g -W -Wall -Wextra

uds_server: uds_server.c
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

.PHONY: clean
clean:
	-rm -f uds_server *.o *~ tmp*
