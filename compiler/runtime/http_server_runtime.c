#include "http_shared.h"

// Forward declaration of fiber functions
extern int64_t fiber_spawn(int64_t (*fn)(void));

// Forward declarations for bridge functions
extern char *process_compile_request(char *json_body);
extern char *process_run_request(char *json_body);

// Thread-safe server context structure
typedef struct {
  int64_t server_id;
} ServerContext;

// Simple HTTP response
static const char *simple_response = "HTTP/1.1 200 OK\r\n"
                                     "Content-Type: text/plain\r\n"
                                     "Content-Length: 13\r\n"
                                     "Connection: close\r\n"
                                     "\r\n"
                                     "Hello, World!";

// Server loop fiber function that actually handles requests
static int64_t server_loop_fiber(void) {
  // Get server from runtime (thread-safe)
  pthread_mutex_lock(&runtime_mutex);
  HttpServer *server = NULL;

  // Find the listening server (there should only be one at a time)
  for (int i = 1; i < MAX_SERVERS; i++) {
    if (servers[i] && servers[i]->is_listening) {
      server = servers[i];
      break;
    }
  }
  pthread_mutex_unlock(&runtime_mutex);

  if (!server) {
    return -1;
  }

  // Keep accepting connections in a loop
  while (server->is_listening) {
    struct sockaddr_in client_addr;
    socklen_t client_len = sizeof(client_addr);

    int client_fd =
        accept(server->socket_fd, (struct sockaddr *)&client_addr, &client_len);
    if (client_fd >= 0) {
      // Read the full HTTP request
      char buffer[4096];
      ssize_t bytes_read = recv(client_fd, buffer, sizeof(buffer) - 1, 0);
      if (bytes_read > 0) {
        buffer[bytes_read] = '\0';

        // Parse the request line to get method and path
        char method[16] = {0};
        char path[256] = {0};
        sscanf(buffer, "%15s %255s", method, path);

        printf("🌐 HTTP Request: %s %s\n", method, path);

        // Find the request body (after \r\n\r\n)
        char *body_start = strstr(buffer, "\r\n\r\n");
        char *body = "";
        if (body_start) {
          body = body_start + 4;
        }

        // Route the request based on path and method
        char response[4096];
        if (strcmp(method, "POST") == 0 && strcmp(path, "/api/compile") == 0) {
          // Call Osprey compile function
          printf("📝 Compile request: %s\n", body);
          char *result = process_compile_request(body);
          if (result) {
            snprintf(response, sizeof(response),
                     "HTTP/1.1 200 OK\r\n"
                     "Content-Type: application/json\r\n"
                     "Content-Length: %zu\r\n"
                     "Connection: close\r\n"
                     "\r\n%s",
                     strlen(result), result);
            free(result);
          } else {
            strcpy(response, "HTTP/1.1 500 Internal Server Error\r\n\r\nError "
                             "processing request");
          }
        } else if (strcmp(method, "POST") == 0 &&
                   strcmp(path, "/api/run") == 0) {
          // Call Osprey run function
          printf("🚀 Run request: %s\n", body);
          char *result = process_run_request(body);
          if (result) {
            snprintf(response, sizeof(response),
                     "HTTP/1.1 200 OK\r\n"
                     "Content-Type: application/json\r\n"
                     "Content-Length: %zu\r\n"
                     "Connection: close\r\n"
                     "\r\n%s",
                     strlen(result), result);
            free(result);
          } else {
            strcpy(response, "HTTP/1.1 500 Internal Server Error\r\n\r\nError "
                             "processing request");
          }
        } else {
          // Default response for other paths
          strcpy(response, simple_response);
        }

        // Send the response
        send(client_fd, response, strlen(response), 0);
      }
      close(client_fd);
    }
  }

  return 0;
}

// Create HTTP server - returns server_id or negative error
int64_t http_create_server(int64_t port, char *address) {
  if (port < 1 || port > 65535) {
    return -1;
  }

  if (!address) {
    return -2;
  }

  int64_t id = get_next_id();
  HttpServer *server = malloc(sizeof(HttpServer));
  if (!server) {
    return -3;
  }

  server->id = id;
  server->port = (int)port;
  server->address = strdup(address);
  if (!server->address) {
    free(server);
    return -3;
  }
  server->socket_fd = -1;
  server->is_listening = false;
  pthread_mutex_init(&server->mutex, NULL);

  pthread_mutex_lock(&runtime_mutex);
  servers[id] = server;
  pthread_mutex_unlock(&runtime_mutex);

  return id;
}

// Start HTTP server listening - returns 0 on success
int64_t http_listen(int64_t server_id, int64_t handler) {
  pthread_mutex_lock(&runtime_mutex);
  HttpServer *server = servers[server_id];
  pthread_mutex_unlock(&runtime_mutex);

  if (!server) {
    return -1;
  }

  // Create socket
  server->socket_fd = socket(AF_INET, SOCK_STREAM, 0);
  if (server->socket_fd < 0) {
    return -2;
  }

  // Set socket options
  int opt = 1;
  if (setsockopt(server->socket_fd, SOL_SOCKET, SO_REUSEADDR, &opt,
                 sizeof(opt)) < 0) {
    close(server->socket_fd);
    return -3;
  }

  // Bind socket
  struct sockaddr_in server_addr;
  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(server->port);
  server_addr.sin_addr.s_addr = inet_addr(server->address);

  if (bind(server->socket_fd, (struct sockaddr *)&server_addr,
           sizeof(server_addr)) < 0) {
    close(server->socket_fd);
    return -4;
  }

  // Start listening
  if (listen(server->socket_fd, SOMAXCONN) < 0) {
    close(server->socket_fd);
    return -5;
  }

  server->is_listening = true;

  // Spawn a fiber to handle the server loop (non-blocking)
  int64_t fiber_id = fiber_spawn(server_loop_fiber);
  if (fiber_id < 0) {
    server->is_listening = false;
    close(server->socket_fd);
    return -6;
  }

  printf("HTTP server listening on %s:%d\n", server->address, server->port);

  return 0;
}

// Stop HTTP server - returns 0 on success
int64_t http_stop_server(int64_t server_id) {
  pthread_mutex_lock(&runtime_mutex);
  HttpServer *server = servers[server_id];
  if (server) {
    servers[server_id] = NULL;
    server->is_listening = false;
    if (server->socket_fd >= 0) {
      close(server->socket_fd);
    }
    free(server->address);
    pthread_mutex_destroy(&server->mutex);
    free(server);
  }
  pthread_mutex_unlock(&runtime_mutex);

  return 0;
}