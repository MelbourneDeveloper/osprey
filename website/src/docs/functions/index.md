---
layout: page
title: "Built-in Functions"
description: "Complete reference for all built-in functions in Osprey"
---

All built-in functions available in Osprey.

## [Channel](channel/)

**Signature:** `Channel(capacity: int) -> Channel`

Creates a new channel with the specified capacity.

## [List](list/)

**Signature:** `List() -> List<T>`

Creates a new empty list.

## [Map](map/)

**Signature:** `Map() -> Map<K, V>`

Creates a new empty map.

## [awaitProcess](awaitprocess/)

**Signature:** `awaitProcess(handle: int) -> int`

Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes.

## [cleanupProcess](cleanupprocess/)

**Signature:** `cleanupProcess(handle: int) -> Unit`

Cleans up resources associated with a completed process. Should be called after awaitProcess.

## [contains](contains/)

**Signature:** `contains(s: string, needle: string) -> bool`

True if needle appears anywhere in s. Empty needle returns true.

## [drop](drop/)

**Signature:** `drop(s: string, n: int) -> string`

Returns s without its first n bytes. Clamps; never fails.

## [endsWith](endswith/)

**Signature:** `endsWith(s: string, suffix: string) -> bool`

True if s ends with suffix.

## [fiber_await](fiber_await/)

**Signature:** `fiber_await(fiber: Fiber) -> any`

Waits for a fiber to complete and returns its result.

## [fiber_spawn](fiber_spawn/)

**Signature:** `fiber_spawn(fn: () -> any) -> Fiber`

Spawns a new fiber to execute the given function concurrently.

## [fiber_yield](fiber_yield/)

**Signature:** `fiber_yield(value: any) -> any`

Yields control to the fiber scheduler with an optional value.

## [filter](filter/)

**Signature:** `filter(iterator: iterator, predicate: function) -> iterator`

Filters elements in an iterator based on a predicate function.

## [fold](fold/)

**Signature:** `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`

Reduces an iterator to a single value by repeatedly applying a function.

## [forEach](foreach/)

**Signature:** `forEach(iterator: iterator, function: function) -> int`

Applies a function to each element in an iterator.

## [forEachList](foreachlist/)

**Signature:** `forEachList(list: List<T>, function: fn(T) -> Unit) -> List<T>`

Apply function to every element of list. Phase 7 of collections plan.

## [httpCloseClient](httpcloseclient/)

**Signature:** `httpCloseClient(clientID: int) -> int`

Closes the HTTP client and cleans up resources.

## [httpCreateClient](httpcreateclient/)

**Signature:** `httpCreateClient(base_url: string, timeout: int) -> int`

Creates an HTTP client for making requests to a base URL.

## [httpCreateServer](httpcreateserver/)

**Signature:** `httpCreateServer(port: int, address: string) -> int`

Creates an HTTP server bound to the specified port and address.

## [httpDelete](httpdelete/)

**Signature:** `httpDelete(clientID: int, path: string, headers: string) -> int`

Makes an HTTP DELETE request to the specified path.

## [httpGet](httpget/)

**Signature:** `httpGet(clientID: int, path: string, headers: string) -> int`

Makes an HTTP GET request to the specified path.

## [httpListen](httplisten/)

**Signature:** `httpListen(serverID: int, handler: (string, string, string, string) -> HttpResponse) -> int`

Starts the HTTP server listening for requests with a handler function.

## [httpPost](httppost/)

**Signature:** `httpPost(clientID: int, path: string, body: string, headers: string) -> int`

Makes an HTTP POST request with a request body.

## [httpPut](httpput/)

**Signature:** `httpPut(clientID: int, path: string, body: string, headers: string) -> int`

Makes an HTTP PUT request with a request body.

## [httpRequest](httprequest/)

**Signature:** `httpRequest(clientID: int, method: int, path: string, headers: string, body: string) -> int`

Makes a generic HTTP request with any method.

## [httpStopServer](httpstopserver/)

**Signature:** `httpStopServer(serverID: int) -> int`

Stops the HTTP server and closes all connections.

## [indexOf](indexof/)

**Signature:** `indexOf(s: string, needle: string) -> Result<int, StringError>`

Returns byte-index of first occurrence of needle, or Error(NotFound).

## [input](input/)

**Signature:** `input() -> Result<string, Error>`

Reads a string from the user's input.

## [isEmpty](isempty/)

**Signature:** `isEmpty(s: string) -> bool`

True if string has zero length.

## [join](join/)

**Signature:** `join(parts: List<string>, separator: string) -> string`

Concatenates parts with separator between each pair.

## [length](length/)

**Signature:** `length(s: string) -> int`

Returns the byte length of a string. Total — never fails.

## [lines](lines/)

**Signature:** `lines(s: string) -> List<string>`

Splits on '\n'. A trailing newline does not produce an empty entry.

## [listAppend](listappend/)

**Signature:** `listAppend(list: List<T>, value: T) -> List<T>`

Returns a new list with value at the end. O(log32 n) amortised.

## [listConcat](listconcat/)

**Signature:** `listConcat(left: List<T>, right: List<T>) -> List<T>`

Returns left ++ right. Same as left + right.

## [listContains](listcontains/)

**Signature:** `listContains(list: List<T>, value: T) -> bool`

True iff some element equals value. O(n).

## [listLength](listlength/)

**Signature:** `listLength(list: List<T>) -> int`

Returns the number of elements in a list. O(1).

## [listPrepend](listprepend/)

**Signature:** `listPrepend(list: List<T>, value: T) -> List<T>`

Returns a new list with value at the front. O(n).

## [listReverse](listreverse/)

**Signature:** `listReverse(list: List<T>) -> List<T>`

Returns a new list in reverse order.

## [map](map/)

**Signature:** `map(iterator: iterator, fn: function) -> iterator`

Transforms each element in an iterator using a function, returning a new iterator.

## [mapContains](mapcontains/)

**Signature:** `mapContains(map: Map<K, V>, key: K) -> bool`

True iff key is present in map.

## [mapKeys](mapkeys/)

**Signature:** `mapKeys(map: Map<K, V>) -> List<K>`

All keys of the map as a list. Order unspecified.

## [mapLength](maplength/)

**Signature:** `mapLength(map: Map<K, V>) -> int`

Returns the number of entries in a map. O(1).

## [mapMerge](mapmerge/)

**Signature:** `mapMerge(left: Map<K, V>, right: Map<K, V>) -> Map<K, V>`

Right-biased union. Same as left + right.

## [mapRemove](mapremove/)

**Signature:** `mapRemove(map: Map<K, V>, key: K) -> Map<K, V>`

Returns a new map without key. No-op if key is absent.

## [mapSet](mapset/)

**Signature:** `mapSet(map: Map<K, V>, key: K, value: V) -> Map<K, V>`

Returns a new map with key bound to value (replaces prior binding).

## [mapValues](mapvalues/)

**Signature:** `mapValues(map: Map<K, V>) -> List<V>`

All values of the map as a list. Order matches mapKeys.

## [padEnd](padend/)

**Signature:** `padEnd(s: string, targetLength: int, fill: string) -> Result<string, StringError>`

Pads s on the right with copies of fill to reach targetLength bytes.

## [padStart](padstart/)

**Signature:** `padStart(s: string, targetLength: int, fill: string) -> Result<string, StringError>`

Pads s on the left with copies of fill to reach targetLength bytes.

## [parseFloat](parsefloat/)

**Signature:** `parseFloat(s: string) -> Result<float, StringError>`

Strict base-10 floating-point parser. No whitespace tolerance.

## [parseInt](parseint/)

**Signature:** `parseInt(s: string) -> Result<int, StringError>`

Strict base-10 signed-int parser. No whitespace tolerance.

## [print](print/)

**Signature:** `print(value: any) -> Unit`

Prints a value to the console. Automatically converts the value to a string representation.

## [range](range/)

**Signature:** `range(start: int, end: int) -> iterator`

Creates an iterator that generates numbers from start to end (exclusive).

## [readFile](readfile/)

**Signature:** `readFile(filename: string) -> Result<string, Error>`

Reads the entire contents of a file as a string.

## [recv](recv/)

**Signature:** `recv(channel: Channel) -> any`

Receives a value from a channel.

## [repeat](repeat/)

**Signature:** `repeat(s: string, n: int) -> Result<string, StringError>`

Concatenates s with itself n times. Error(InvalidArgument) on negative n.

## [replace](replace/)

**Signature:** `replace(s: string, needle: string, replacement: string) -> Result<string, StringError>`

Replaces every occurrence of needle. Error(InvalidArgument) on empty needle.

## [reverse](reverse/)

**Signature:** `reverse(s: string) -> string`

Reverses byte order. Grapheme-cluster reversal is future work.

## [send](send/)

**Signature:** `send(channel: Channel, value: any) -> int`

Sends a value to a channel. Returns 1 for success, 0 for failure.

## [sleep](sleep/)

**Signature:** `sleep(milliseconds: int) -> int`

Pauses execution for the specified number of milliseconds.

## [spawnProcess](spawnprocess/)

**Signature:** `spawnProcess(command: string, callback: (int, int, string) -> Unit) -> Result<ProcessHandle, string>`

Spawns an external async process with MANDATORY callback for stdout/stderr capture. The callback function receives (processID: int, eventType: int, data: string) and is called for stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING!

## [split](split/)

**Signature:** `split(s: string, separator: string) -> Result<List<string>, StringError>`

Splits s on separator. Error(InvalidArgument) on empty separator.

## [startsWith](startswith/)

**Signature:** `startsWith(s: string, prefix: string) -> bool`

True if s begins with prefix.

## [substring](substring/)

**Signature:** `substring(s: string, start: int, end: int) -> Result<string, StringError>`

Extracts s[start, end). Returns Error(IndexOutOfRange) if start<0, end>len, or start>end.

## [take](take/)

**Signature:** `take(s: string, n: int) -> string`

Returns at most the first n bytes of s. Clamps; never fails.

## [toLowerCase](tolowercase/)

**Signature:** `toLowerCase(s: string) -> string`

ASCII-aware lowercase.

## [toString](tostring/)

**Signature:** `toString(value: any) -> string`

Converts a value to its string representation.

## [toUpperCase](touppercase/)

**Signature:** `toUpperCase(s: string) -> string`

ASCII-aware uppercase. Unicode simple case mapping is a future addition.

## [trim](trim/)

**Signature:** `trim(s: string) -> string`

Removes leading and trailing whitespace.

## [trimEnd](trimend/)

**Signature:** `trimEnd(s: string) -> string`

Removes trailing whitespace.

## [trimStart](trimstart/)

**Signature:** `trimStart(s: string) -> string`

Removes leading whitespace.

## [websocketClose](websocketclose/)

**Signature:** `websocketClose(wsID: Int) -> Result<Success, String>`

Closes the WebSocket connection and cleans up resources.

## [websocketConnect](websocketconnect/)

**Signature:** `websocketConnect(url: String, messageHandler: (String) -> Result<Success, String>) -> Result<WebSocketID, String>`

Establishes a WebSocket connection with a message handler callback.

## [websocketCreateServer](websocketcreateserver/)

**Signature:** `websocketCreateServer(port: Int, address: String, path: String) -> Result<ServerID, String>`

Creates a WebSocket server bound to the specified port, address, and path.

## [websocketKeepAlive](websocketkeepalive/)

**Signature:** `websocketKeepAlive() -> Unit`

Keeps the WebSocket server running indefinitely until interrupted (blocking operation).

## [websocketSend](websocketsend/)

**Signature:** `websocketSend(wsID: Int, message: String) -> Result<Success, String>`

Sends a message through the WebSocket connection.

## [websocketServerBroadcast](websocketserverbroadcast/)

**Signature:** `websocketServerBroadcast(serverID: Int, message: String) -> Result<Success, String>`

Broadcasts a message to all connected WebSocket clients.

## [websocketServerListen](websocketserverlisten/)

**Signature:** `websocketServerListen(serverID: Int) -> Result<Success, String>`

Starts the WebSocket server listening for connections.

## [websocketStopServer](websocketstopserver/)

**Signature:** `websocketStopServer(serverID: Int) -> Result<Success, String>`

Stops the WebSocket server and closes all connections.

## [words](words/)

**Signature:** `words(s: string) -> List<string>`

Splits on runs of whitespace; empty results dropped.

## [writeFile](writefile/)

**Signature:** `writeFile(filename: string, content: string) -> Result<int, Error>`

Writes content to a file. Creates the file if it doesn't exist. Returns number of bytes written.

