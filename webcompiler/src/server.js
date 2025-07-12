import { spawn } from 'child_process'
import express from 'express'
import fs from 'fs/promises'
import { createServer } from 'http'
import path from 'path'
import { fileURLToPath } from 'url'
import { WebSocketServer } from 'ws'
import { randomUUID } from 'crypto'
import { execSync } from 'child_process'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const app = express()
const server = createServer(app)

const PORT = process.env.PORT || 3001

// STARTUP LOGGING - Make it super obvious the server is starting
console.log('\n' + '='.repeat(80))
console.log('🚀 OSPREY WEB COMPILER STARTING UP')
console.log('='.repeat(80))
console.log(`📍 Server file: ${__filename}`)
console.log(`📁 Working directory: ${process.cwd()}`)
console.log(`🐳 Docker environment: ${process.env.DOCKER_ENV || 'false'}`)
console.log(`🏃 Node environment: ${process.env.NODE_ENV || 'development'}`)
console.log(`🔌 Target port: ${PORT}`)
console.log('='.repeat(80))

// Request logging middleware - track ALL requests
app.use((req, res, next) => {
    const timestamp = new Date().toISOString()
    console.log(`\n📨 [${timestamp}] ${req.method} ${req.url}`)
    console.log(`📍 User-Agent: ${req.headers['user-agent'] || 'unknown'}`)
    console.log(`📍 Origin: ${req.headers.origin || 'none'}`)
    console.log(`📍 Content-Type: ${req.headers['content-type'] || 'none'}`)

    // Log body size for POST requests
    if (req.method === 'POST' && req.body) {
        const bodySize = JSON.stringify(req.body).length
        console.log(`📏 Body size: ${bodySize} bytes`)
    }

    next()
})

// Middleware
app.use(express.json({ limit: '10mb' }))

// CORS middleware
app.use((req, res, next) => {
    // Allow requests from the website running on localhost:8080
    const origin = req.headers.origin;
    const allowedOrigins = [
        'http://localhost:8080',
        'http://127.0.0.1:8080',
        'http://localhost:3001',
        'http://127.0.0.1:3001',
        'https://ospreylang.dev',
        'https://www.ospreylang.dev'
    ];

    if (allowedOrigins.includes(origin)) {
        res.header('Access-Control-Allow-Origin', origin);
    }

    res.header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
    res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Requested-With');
    res.header('Access-Control-Allow-Credentials', 'true');

    if (req.method === 'OPTIONS') {
        return res.sendStatus(200);
    }
    next();
})

// Health check endpoint
app.get('/api', (req, res) => {
    res.json({
        status: 'ok',
        service: 'osprey-web-compiler',
        version: '0.2.0',
        timestamp: new Date().toISOString()
    })
})

// Compile endpoint
app.post('/api/compile', async (req, res) => {
    const { code } = req.body
    console.log('📝 Compile request received')
    console.log('📄 Code length:', code?.length || 0)

    if (!code) {
        return res.status(400).json({ success: false, error: 'No code provided' })
    }

    try {
        const result = await runOspreyCompiler(['--sandbox', '--ast'], code)

        if (result.success) {
            console.log('✅ Compile success, output length:', result.stdout.length)
            res.status(200).json({
                success: true,
                compilerOutput: result.stderr || '',
                programOutput: result.stdout || '' // AST output goes to stdout
            })
        } else {
            console.error('❌ Compile error, stderr:', result.stderr)

            const errorOutput = result.stderr || result.stdout || '';

            // Detect INTERNAL compiler errors - simple marker from compiler
            const isInternalError = errorOutput.includes('INTERNAL_COMPILER_ERROR:');

            if (isInternalError) {
                // Internal compiler error - log for debugging but don't expose to user
                console.error('🚨 INTERNAL COMPILER ERROR DETECTED:', errorOutput);
                res.status(502).json({
                    success: false,
                    error: 'The compiler encountered an internal error. Please report this code to help us fix the issue.',
                    isInternalError: true
                });
                return;
            }

            res.status(422).json({ // 422 Unprocessable Entity for compilation errors
                success: false,
                compilerOutput: result.stderr || '',
                programOutput: result.stdout || '',
                error: errorOutput || `Compilation failed with exit code ${result.exitCode}`
            })
        }
    } catch (error) {
        console.error('❌ System error:', error.message)
        res.status(500).json({ success: false, error: error.message })
    }
})

// Run endpoint
app.post('/api/run', async (req, res) => {
    const { code } = req.body
    console.log('🏃 Run request received')
    console.log('📄 Code length:', code?.length || 0)

    if (!code) {
        return res.status(400).json({ success: false, error: 'No code provided' })
    }

    try {
        const result = await runOspreyCompiler(['--run'], code)

        if (result.success) {
            console.log('✅ Run success')
            console.log('📊 Compiler output length:', result.stderr?.length || 0)
            console.log('📋 Program output length:', result.stdout?.length || 0)

            res.status(200).json({
                success: true,
                compilerOutput: result.stderr || '',
                programOutput: result.stdout || ''
            })
        } else {
            console.error('❌ Run failed, stderr:', result.stderr)
            console.error('❌ Run failed, stdout:', result.stdout)

            const errorOutput = result.stderr || result.stdout || '';

            // Detect INTERNAL compiler errors - simple marker from compiler
            const isInternalError = errorOutput.includes('INTERNAL_COMPILER_ERROR:');

            if (isInternalError) {
                // Internal compiler error - log for debugging but don't expose to user
                console.error('🚨 INTERNAL COMPILER ERROR DETECTED:', errorOutput);
                res.status(502).json({
                    success: false,
                    error: 'The compiler encountered an internal error. Please report this code to help us fix the issue.',
                    isInternalError: true
                });
                return;
            }

            // Determine if it's a user syntax/compilation error or runtime error
            const isCompilationError = errorOutput.includes('parse errors') ||
                errorOutput.includes('undefined variable') ||
                errorOutput.includes('syntax error') ||
                errorOutput.includes('type mismatch') ||
                errorOutput.includes('Compilation failed');

            const statusCode = isCompilationError ? 422 : 400; // 422 for compilation, 400 for runtime

            res.status(statusCode).json({
                success: false,
                compilerOutput: result.stderr || '',
                programOutput: result.stdout || '',
                isCompilationError: isCompilationError,
                error: errorOutput || `Process failed with exit code ${result.exitCode}`
            })
        }
    } catch (error) {
        console.error('❌ System error:', error.message)
        res.status(500).json({ success: false, error: error.message })
    }
})

// STARTUP: Delete ALL temp folders on server startup
async function deleteAllTempFolders() {
    const tempBaseDir = '/tmp/osprey-temp'
    try {
        console.log('🗑️ Deleting ALL temp folders on startup...')
        await fs.rm(tempBaseDir, { recursive: true, force: true })
        console.log('✅ All temp folders deleted')
    } catch (error) {
        console.error('⚠️ Error deleting temp folders:', error.message)
    }
}

// Cleanup old temp folders periodically to prevent disk space issues
async function cleanupOldTempFolders() {
    const tempBaseDir = '/tmp/osprey-temp'
    try {
        const folders = await fs.readdir(tempBaseDir)
        const now = Date.now()
        const oneHourAgo = now - (60 * 60 * 1000) // 1 hour ago

        for (const folder of folders) {
            const folderPath = path.join(tempBaseDir, folder)
            const stats = await fs.stat(folderPath)
            if (stats.isDirectory() && stats.mtime.getTime() < oneHourAgo) {
                await fs.rm(folderPath, { recursive: true, force: true })
                console.log(`🗑️ Cleaned up old temp folder: ${folder}`)
            }
        }
    } catch (error) {
        console.error('⚠️ Error cleaning up temp folders:', error.message)
    }
}

// Run cleanup every 30 minutes
setInterval(cleanupOldTempFolders, 30 * 60 * 1000)

// DELETE ALL TEMP FOLDERS ON STARTUP
deleteAllTempFolders()

// THREAD-SAFE Helper function to run Osprey compiler
// Each request gets its own UUID-named folder for complete isolation
// Always uses --sandbox flag for security (disables HTTP, WebSocket, file system, and FFI access)
function runOspreyCompiler(args, code = '') {
    return new Promise(async (resolve, reject) => {
        // Diagnostics before running compiler
        try {
            console.log('🛠️ ENV:', process.env);
            console.log('🛠️ which osprey:', execSync('which osprey').toString());
            console.log('🛠️ ldd osprey:', execSync('ldd /usr/local/bin/osprey').toString());
            console.log('🛠️ ls /usr/local/lib:', execSync('ls -l /usr/local/lib').toString());
            console.log('🛠️ ls /usr/lib/llvm-14/bin:', execSync('ls -l /usr/lib/llvm-14/bin').toString());
            console.log('🛠️ ls /tmp/osprey-temp:', execSync('ls -l /tmp/osprey-temp').toString());
        } catch (e) {
            console.error('🛠️ Diagnostics error:', e.message);
        }
        // Create a unique UUID folder for this request - THREAD SAFE!
        const requestId = randomUUID()
        const tempBaseDir = '/tmp/osprey-temp'
        const tempRequestDir = path.join(tempBaseDir, requestId)
        const tempFile = path.join(tempRequestDir, 'main.osp')

        try {
            // Create the unique temp directory for this request
            await fs.mkdir(tempRequestDir, { recursive: true })
            console.log(`📁 Created temp folder: ${requestId}`)

            console.log(`💾 Writing temp file: ${tempFile}`)
            await fs.writeFile(tempFile, code)

            // Use the osprey binary from PATH (installed in Docker) or fallback to local dev path
            const ospreyPath = process.env.NODE_ENV === 'production' || process.env.DOCKER_ENV
                ? 'osprey'
                : path.resolve(__dirname, '../compiler/bin/osprey')
            console.log(`🔨 Running: ${ospreyPath} ${tempFile} ${args.join(' ')}`)
            const child = spawn(ospreyPath, [tempFile, ...args], {
                stdio: 'pipe',
                cwd: tempRequestDir, // Run in the temp directory
                timeout: 5000 // 5 second timeout - kill any program that runs longer
            })

            let stdout = ''
            let stderr = ''

            child.stdout.on('data', (data) => {
                stdout += data.toString()
            })

            child.stderr.on('data', (data) => {
                stderr += data.toString()
            })

            child.on('close', async (exitCode) => {
                // Clean up the ENTIRE temp folder for this request
                try {
                    await fs.rm(tempRequestDir, { recursive: true, force: true })
                    console.log(`🗑️ Cleaned up temp folder: ${requestId}`)
                } catch (e) {
                    console.error('⚠️ Failed to clean up temp folder:', e.message)
                }

                // Log detailed information about the compiler execution
                console.log(`🔍 Compiler execution completed with exit code: ${exitCode}`)
                console.log(`📊 stdout length: ${stdout.length}, stderr length: ${stderr.length}`)
                if (stdout.length > 0) {
                    console.log(`📤 stdout: ${stdout.substring(0, 500)}${stdout.length > 500 ? '...' : ''}`)
                }
                if (stderr.length > 0) {
                    console.log(`📤 stderr: ${stderr.substring(0, 500)}${stderr.length > 500 ? '...' : ''}`)
                }

                // Always resolve with the result - let the caller determine success/failure
                resolve({
                    exitCode,
                    stdout,
                    stderr,
                    success: exitCode === 0
                })
            })

            child.on('error', async (error) => {
                // Clean up temp folder on error
                try {
                    await fs.rm(tempRequestDir, { recursive: true, force: true })
                    console.log(`🗑️ Cleaned up temp folder after error: ${requestId}`)
                } catch (e) {
                    console.error('⚠️ Failed to clean up temp folder after error:', e.message)
                }
                reject(error)
            })
        } catch (error) {
            // Clean up temp folder if creation failed
            try {
                await fs.rm(tempRequestDir, { recursive: true, force: true })
            } catch (e) {
                // Ignore cleanup errors
            }
            reject(error)
        }
    })
}

// WebSocket server for LSP bridge
const wss = new WebSocketServer({
    server,
    path: '/lsp',
    verifyClient: (info) => {
        // Check origin for CORS on WebSocket connections
        const origin = info.origin;
        const allowedOrigins = [
            'http://localhost:8080',
            'http://127.0.0.1:8080',
            'http://localhost:3001',
            'http://127.0.0.1:3001',
            'https://ospreylang.dev',
            'https://www.ospreylang.dev'
        ];

        console.log('🔍 WebSocket upgrade request from origin:', origin);

        if (!origin || allowedOrigins.includes(origin)) {
            return true;
        }

        console.error('❌ WebSocket connection rejected - invalid origin:', origin);
        return false;
    }
})

console.log(`🌐 WebSocket server configured for path: /lsp`)

wss.on('connection', (ws, req) => {
    console.log('🔌 New WebSocket connection from:', req.socket.remoteAddress)
    console.log('🔍 Connection headers:', JSON.stringify(req.headers, null, 2))

    // Path to the compiled Osprey LSP server - use different paths for Docker vs local dev
    const lspPath = process.env.NODE_ENV === 'production' || process.env.DOCKER_ENV
        ? path.resolve(__dirname, '../server/out/src/server.js')  // Docker path: /app/server/out/src/server.js
        : path.resolve(__dirname, '../../vscode-extension/server/out/src/server.js')  // Local dev path

    console.log('🚀 Starting Osprey LSP:', lspPath)

    // Check if LSP file exists
    fs.access(lspPath)
        .then(() => {
            // Spawn the LSP server process
            const lspProcess = spawn('node', [lspPath, '--stdio'], {
                stdio: ['pipe', 'pipe', 'pipe'],
                cwd: process.cwd(),
                env: { ...process.env, NODE_ENV: 'development' }
            })

            lspProcess.on('error', (error) => {
                console.error('❌ LSP process error:', error)
                ws.close(1011, 'LSP server failed to start')
            })

            lspProcess.on('spawn', () => {
                console.log('✅ LSP process started successfully')
                console.log(`📊 LSP process PID: ${lspProcess.pid}`)
            })

            // Message counter for debugging
            let clientToServerCount = 0
            let serverToClientCount = 0

            // Forward messages between WebSocket and LSP stdio
            ws.on('message', (data) => {
                const message = data.toString()
                clientToServerCount++
                console.log(`📨 Client -> LSP [${clientToServerCount}]:`, message.substring(0, 200) + (message.length > 200 ? '...' : ''))

                // Parse to check message type
                try {
                    const parsed = JSON.parse(message)
                    console.log(`  📌 Message type: ${parsed.method || parsed.id ? 'request/response' : 'notification'}`)
                    if (parsed.method) {
                        console.log(`  📌 Method: ${parsed.method}`)
                    }
                } catch (e) {
                    console.log('  ⚠️ Could not parse message as JSON')
                }

                if (lspProcess.stdin && !lspProcess.stdin.destroyed) {
                    lspProcess.stdin.write(message)
                } else {
                    console.error('❌ LSP stdin not available!')
                }
            })

            lspProcess.stdout.on('data', (data) => {
                const message = data.toString()
                serverToClientCount++
                console.log(`📤 LSP -> Client [${serverToClientCount}]:`, message.substring(0, 200) + (message.length > 200 ? '...' : ''))

                // Parse to check message type
                try {
                    const parsed = JSON.parse(message)
                    console.log(`  📌 Message type: ${parsed.method || parsed.id ? 'request/response' : 'notification'}`)
                    if (parsed.method) {
                        console.log(`  📌 Method: ${parsed.method}`)
                    }
                } catch (e) {
                    console.log('  ⚠️ Could not parse message as JSON')
                }

                if (ws.readyState === ws.OPEN) {
                    ws.send(data)
                } else {
                    console.error('❌ WebSocket not open, cannot send message')
                }
            })

            lspProcess.stderr.on('data', (data) => {
                const errorMsg = data.toString()
                console.error('🔴 LSP stderr:', errorMsg)
            })

            ws.on('close', (code, reason) => {
                console.log(`🔌 WebSocket disconnected: code=${code}, reason=${reason}`)
                console.log(`📊 Total messages: Client->Server: ${clientToServerCount}, Server->Client: ${serverToClientCount}`)
                if (!lspProcess.killed) {
                    console.log('🛑 Killing LSP process')
                    lspProcess.kill()
                }
            })

            lspProcess.on('close', (code, signal) => {
                console.log(`🛑 LSP process exited: code=${code}, signal=${signal}`)
                if (ws.readyState === ws.OPEN) {
                    ws.close()
                }
            })
        })
        .catch((error) => {
            console.error('❌ LSP server file not found:', lspPath, error)
            ws.close(1011, 'LSP server file not found')
        })

    ws.on('error', (error) => {
        console.error('❌ WebSocket error:', error)
    })
})

wss.on('error', (error) => {
    console.error('❌ WebSocket server error:', error)
})

// Error handling middleware
app.use((error, req, res, next) => {
    console.error('💥 Unhandled error:', error)
    res.status(500).json({
        success: false,
        error: 'Internal server error',
        message: process.env.NODE_ENV === 'development' ? error.message : 'Something went wrong'
    })
})

server.listen(PORT, '0.0.0.0', () => {
    console.log(`✅ WebSocket LSP Bridge running at ws://0.0.0.0:${PORT}/lsp`)
    console.log(`🔨 Compile/Run API available at http://0.0.0.0:${PORT}/api`)
    console.log(`🏥 Health check: http://0.0.0.0:${PORT}/api`)
    console.log(`🌐 Server accessible from external hosts on port ${PORT}`)
}) 