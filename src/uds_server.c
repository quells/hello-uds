#include <errno.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/un.h>
#include <unistd.h>

#define BUF_SIZE 4096

const char *SOCK_FILE = "/tmp/hello-uds.sock";

void handle_conn(int cli_fd, char *buf);

int main(int argc, char **argv) {
    int exit_code = EXIT_SUCCESS;
    int sock_fd, cli_fd;
    char *buf;

    if ((sock_fd = socket(PF_UNIX, SOCK_STREAM, 0)) < 0) {
        perror("socket");
        return EXIT_FAILURE;
    }

    struct sockaddr_un srv_addr, cli_addr;
    memset(&srv_addr, 0, sizeof(srv_addr));
    srv_addr.sun_family = AF_UNIX;
    strcpy(srv_addr.sun_path, SOCK_FILE);
    if (unlink(SOCK_FILE) < 0) {
        perror("unlink");
    }

    if (bind(sock_fd, (struct sockaddr *)&srv_addr, sizeof(srv_addr)) < 0) {
        perror("bind");
        exit_code = EXIT_FAILURE;
        goto CLEANUP;
    }

    if (listen(sock_fd, 0) < 0) {
        perror("listen");
        exit_code = EXIT_FAILURE;
        goto CLEANUP;
    }

    buf = malloc(BUF_SIZE * sizeof(char));
    if (buf == NULL) {
        perror("malloc buf");
        exit_code = EXIT_FAILURE;
        goto CLEANUP;
    }

    socklen_t cli_len = sizeof(cli_addr);
    do {
        cli_fd = accept(sock_fd, (struct sockaddr *)&cli_addr, &cli_len);
        if (cli_fd < 0) {
            perror("accept");
            exit_code = EXIT_FAILURE;
            goto CLEANUP;
        }

        handle_conn(cli_fd, buf);
    } while (1);

CLEANUP:
    if (sock_fd >= 0) {
        if (close(sock_fd) < 0) {
            perror("closing sock_fd");
        }
    }
    if (buf != NULL) {
        free(buf);
    }

    return exit_code;
}

void handle_conn(int cli_fd, char *buf) {
    do {
        int n = read(cli_fd, buf, BUF_SIZE - 1);
        if (n < 0) {
            perror("reading from client");
            goto HANDLE_CONN_CLEANUP;
        }
        buf[n] = '\0';
        printf("got %s\n", buf);
        if (write(cli_fd, "OK\n", 3) < 0) {
            perror("writing response");
            goto HANDLE_CONN_CLEANUP;
        }
    } while (1);
HANDLE_CONN_CLEANUP:
    if (close(cli_fd) < 0) {
        perror("closing cli_fd");
    }
}
