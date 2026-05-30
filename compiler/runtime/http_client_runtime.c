#include "http_shared.h"

// HTTP client runtime.
//
// Two request surfaces share one transport helper (http_perform):
//   * Status-only API  (http_get/post/put/delete, http_request) returns the
//     numeric status code, preserving the original behaviour.
//   * Response-handle API (http_*_response + http_response_*) retains the full
//     response body and headers in a heap slot and returns an opaque handle so
//     Osprey can read the body. Implements [HTTP-RESPONSE-HANDLE].

// Create HTTP client - returns client_id or negative error
int64_t http_create_client(char *base_url, int64_t timeout) {
  if (!base_url) {
    return -1;
  }

  if (timeout < 0) {
    return -2;
  }

  int64_t id = get_next_id();
  HttpClient *client = malloc(sizeof(HttpClient));
  if (!client) {
    return -3;
  }

  client->id = id;
  client->base_url = strdup(base_url);
  client->timeout = (int)timeout;
  client->is_persistent = false;

  // Parse base URL
  char *path;
  if (parse_url(base_url, &client->host, &client->port, &path) != 0) {
    free(client->base_url);
    free(client);
    return -4;
  }
  free(path); // We only need host and port for client

  pthread_mutex_lock(&runtime_mutex);
  clients[id] = client;
  pthread_mutex_unlock(&runtime_mutex);

  return id;
}

// http_perform connects, sends the request and reads the ENTIRE response into a
// freshly allocated, NUL-terminated buffer (the server uses Connection: close,
// so recv() draining to EOF captures the whole body). Returns 0 on success with
// *out_raw / *out_len set (caller frees *out_raw), or a negative error code.
static int64_t http_perform(int64_t client_id, int64_t method, char *path,
                            char *headers, char *body, char **out_raw,
                            size_t *out_len) {
  pthread_mutex_lock(&runtime_mutex);
  HttpClient *client = clients[client_id];
  pthread_mutex_unlock(&runtime_mutex);

  if (!client) {
    return -1;
  }
  if (!path) {
    return -2;
  }

  int sock = socket(AF_INET, SOCK_STREAM, 0);
  if (sock < 0) {
    return -3;
  }

  struct timeval tv;
  tv.tv_sec = client->timeout / 1000;
  tv.tv_usec = (client->timeout % 1000) * 1000;
  setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, (const char *)&tv, sizeof tv);

  struct sockaddr_in server_addr;
  memset(&server_addr, 0, sizeof(server_addr));
  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(client->port);

  struct hostent *server = gethostbyname(client->host);
  if (!server) {
#ifdef _WIN32
    closesocket(sock);
#else
    close(sock);
#endif
    return -4;
  }

  memcpy(&server_addr.sin_addr.s_addr, server->h_addr, server->h_length);

  if (connect(sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
#ifdef _WIN32
    closesocket(sock);
#else
    close(sock);
#endif
    return -5;
  }

  char request[MAX_HTTP_BUFFER];
  const char *method_str = http_method_to_string((HttpMethod)method);

  int request_len = snprintf(request, sizeof(request),
                             "%s %s HTTP/1.1\r\n"
                             "Host: %s:%d\r\n"
                             "Connection: close\r\n",
                             method_str, path, client->host, client->port);

  if (headers && strlen(headers) > 0) {
    request_len += snprintf(request + request_len, sizeof(request) - request_len,
                            "%s\r\n", headers);
  }

  if (body && strlen(body) > 0) {
    request_len += snprintf(request + request_len, sizeof(request) - request_len,
                            "Content-Length: %zu\r\n\r\n%s", strlen(body), body);
  } else {
    request_len += snprintf(request + request_len,
                            sizeof(request) - request_len, "\r\n");
  }

  if (send(sock, request, request_len, 0) < 0) {
#ifdef _WIN32
    closesocket(sock);
#else
    close(sock);
#endif
    return -6;
  }

  size_t cap = MAX_HTTP_BUFFER;
  size_t len = 0;
  char *buf = malloc(cap);
  if (!buf) {
#ifdef _WIN32
    closesocket(sock);
#else
    close(sock);
#endif
    return -9;
  }

  for (;;) {
    if (len + CHUNK_SIZE + 1 > cap) {
      cap *= 2;
      char *nb = realloc(buf, cap);
      if (!nb) {
        free(buf);
#ifdef _WIN32
        closesocket(sock);
#else
        close(sock);
#endif
        return -9;
      }
      buf = nb;
    }
    ssize_t n = recv(sock, buf + len, CHUNK_SIZE, 0);
    if (n > 0) {
      len += (size_t)n;
    } else {
      break; // 0 = peer closed, <0 = timeout/error: stop draining
    }
  }

#ifdef _WIN32
  closesocket(sock);
#else
  close(sock);
#endif

  if (len == 0) {
    free(buf);
    return -7;
  }

  buf[len] = '\0';
  *out_raw = buf;
  *out_len = len;
  return 0;
}

// parse_status_line returns the numeric status from an HTTP status line, or -8.
static int64_t parse_status_line(const char *raw) {
  if (strncmp(raw, "HTTP/1.1 ", 9) == 0 || strncmp(raw, "HTTP/1.0 ", 9) == 0) {
    return atoi(raw + 9);
  }
  return -8;
}

// Make HTTP request - returns HTTP status code or negative error
int64_t http_request(int64_t client_id, int64_t method, char *path,
                     char *headers, char *body) {
  char *raw = NULL;
  size_t len = 0;
  int64_t rc = http_perform(client_id, method, path, headers, body, &raw, &len);
  if (rc != 0) {
    return rc;
  }
  int64_t status = parse_status_line(raw);
  free(raw);
  return status;
}

// Close HTTP client - returns 0 on success
int64_t http_close_client(int64_t client_id) {
  pthread_mutex_lock(&runtime_mutex);
  HttpClient *client = clients[client_id];
  if (client) {
    clients[client_id] = NULL;
    free(client->base_url);
    free(client->host);
    free(client);
  }
  pthread_mutex_unlock(&runtime_mutex);

  return 0;
}

// Convenience functions for specific HTTP methods (status-only API)
int64_t http_get(int64_t client_id, char *path, char *headers) {
  return http_request(client_id, HTTP_GET, path, headers, NULL);
}

int64_t http_post(int64_t client_id, char *path, char *body, char *headers) {
  return http_request(client_id, HTTP_POST, path, headers, body);
}

int64_t http_put(int64_t client_id, char *path, char *body, char *headers) {
  return http_request(client_id, HTTP_PUT, path, headers, body);
}

int64_t http_delete(int64_t client_id, char *path, char *headers) {
  return http_request(client_id, HTTP_DELETE, path, headers, NULL);
}

// ---------------------------------------------------------------------------
// Response-handle API. Implements [HTTP-RESPONSE-HANDLE].
// ---------------------------------------------------------------------------

#define MAX_HTTP_RESPONSES 1024

typedef struct {
  bool in_use;
  int64_t status;
  char *body;        // response body (NUL-terminated, heap)
  char *raw_headers; // header block between status line and blank line (heap)
} HttpClientResponse;

static HttpClientResponse g_http_responses[MAX_HTTP_RESPONSES];
static pthread_mutex_t g_http_resp_mutex = PTHREAD_MUTEX_INITIALIZER;

static bool valid_response_handle(int64_t handle) {
  return handle >= 1 && handle < MAX_HTTP_RESPONSES;
}

// http_request_capture performs the request, retains status/body/headers in a
// heap slot, and returns a 1-based handle (>=1) or a negative error code.
int64_t http_request_capture(int64_t client_id, int64_t method, char *path,
                             char *headers, char *body) {
  char *raw = NULL;
  size_t len = 0;
  int64_t rc = http_perform(client_id, method, path, headers, body, &raw, &len);
  if (rc != 0) {
    return rc;
  }

  int64_t status = parse_status_line(raw);

  char *sep = strstr(raw, "\r\n\r\n");
  char *body_start = sep ? sep + 4 : raw;
  char *first_eol = strstr(raw, "\r\n");

  char *resp_headers = NULL;
  if (first_eol && sep && sep >= first_eol + 2) {
    size_t hlen = (size_t)(sep - (first_eol + 2));
    resp_headers = malloc(hlen + 1);
    if (resp_headers) {
      memcpy(resp_headers, first_eol + 2, hlen);
      resp_headers[hlen] = '\0';
    }
  }

  char *resp_body = strdup(body_start);
  free(raw);

  if (!resp_body) {
    free(resp_headers);
    return -9;
  }

  pthread_mutex_lock(&g_http_resp_mutex);
  int64_t handle = -1;
  for (int64_t i = 1; i < MAX_HTTP_RESPONSES; i++) {
    if (!g_http_responses[i].in_use) {
      g_http_responses[i].in_use = true;
      g_http_responses[i].status = status;
      g_http_responses[i].body = resp_body;
      g_http_responses[i].raw_headers = resp_headers;
      handle = i;
      break;
    }
  }
  pthread_mutex_unlock(&g_http_resp_mutex);

  if (handle < 0) {
    free(resp_body);
    free(resp_headers);
    return -10; // response table full
  }
  return handle;
}

int64_t http_get_response(int64_t client_id, char *path, char *headers) {
  return http_request_capture(client_id, HTTP_GET, path, headers, NULL);
}

int64_t http_post_response(int64_t client_id, char *path, char *body,
                           char *headers) {
  return http_request_capture(client_id, HTTP_POST, path, headers, body);
}

int64_t http_response_status(int64_t handle) {
  if (!valid_response_handle(handle)) {
    return -1;
  }
  pthread_mutex_lock(&g_http_resp_mutex);
  int64_t status = g_http_responses[handle].in_use
                       ? g_http_responses[handle].status
                       : -1;
  pthread_mutex_unlock(&g_http_resp_mutex);
  return status;
}

// Returns a fresh copy of the body so it survives http_response_free.
char *http_response_body(int64_t handle) {
  if (!valid_response_handle(handle)) {
    return NULL;
  }
  pthread_mutex_lock(&g_http_resp_mutex);
  char *copy = NULL;
  if (g_http_responses[handle].in_use && g_http_responses[handle].body) {
    copy = strdup(g_http_responses[handle].body);
  }
  pthread_mutex_unlock(&g_http_resp_mutex);
  return copy;
}

// Case-insensitive prefix match of "name:" at the start of a header line.
static char *find_header_value(const char *headers, const char *name) {
  if (!headers || !name) {
    return NULL;
  }
  size_t name_len = strlen(name);
  const char *line = headers;
  while (line && *line) {
    bool match = true;
    for (size_t i = 0; i < name_len; i++) {
      char a = line[i];
      char b = name[i];
      if (a >= 'A' && a <= 'Z') {
        a = (char)(a - 'A' + 'a');
      }
      if (b >= 'A' && b <= 'Z') {
        b = (char)(b - 'A' + 'a');
      }
      if (a != b || a == '\0') {
        match = false;
        break;
      }
    }
    if (match && line[name_len] == ':') {
      const char *value = line + name_len + 1;
      while (*value == ' ' || *value == '\t') {
        value++;
      }
      const char *end = strstr(value, "\r\n");
      size_t vlen = end ? (size_t)(end - value) : strlen(value);
      char *out = malloc(vlen + 1);
      if (!out) {
        return NULL;
      }
      memcpy(out, value, vlen);
      out[vlen] = '\0';
      return out;
    }
    const char *next = strstr(line, "\r\n");
    line = next ? next + 2 : NULL;
  }
  return NULL;
}

char *http_response_header(int64_t handle, char *name) {
  if (!valid_response_handle(handle)) {
    return NULL;
  }
  pthread_mutex_lock(&g_http_resp_mutex);
  char *value = NULL;
  if (g_http_responses[handle].in_use) {
    value = find_header_value(g_http_responses[handle].raw_headers, name);
  }
  pthread_mutex_unlock(&g_http_resp_mutex);
  return value;
}

// Frees a response slot. Returns 0 on success, -1 on invalid/double free.
int64_t http_response_free(int64_t handle) {
  if (!valid_response_handle(handle)) {
    return -1;
  }
  pthread_mutex_lock(&g_http_resp_mutex);
  int64_t rc = -1;
  if (g_http_responses[handle].in_use) {
    free(g_http_responses[handle].body);
    free(g_http_responses[handle].raw_headers);
    g_http_responses[handle].body = NULL;
    g_http_responses[handle].raw_headers = NULL;
    g_http_responses[handle].in_use = false;
    rc = 0;
  }
  pthread_mutex_unlock(&g_http_resp_mutex);
  return rc;
}
