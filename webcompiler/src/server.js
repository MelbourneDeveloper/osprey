import { spawn } from 'child_process'
import { randomUUID } from 'crypto'
import express from 'express'
import fs from 'fs/promises'
import { createServer } from 'http'
import path from 'path'
import { fileURLToPath } from 'url'
import { WebSocketServer } from 'ws'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const app = express()
const server = createServer(app)

const PORT = process.env.PORT || 3001

// Set umask to ensure files are created with execute permissions
process.umask(0o000)

console.log('🚀 Starting WebSocket LSP Bridge...')

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
    console.log('📋 COMPILE REQUEST:')
    console.log(code || '(no code)')

    if (!code) {
        return res.status(400).json({ success: false, error: 'No code provided' })
    }

    try {
        const result = await runOspreyCompiler(['--sandbox', '--ast'], code)

        if (result.success) {
            res.status(200).json({
                success: true,
                compilerOutput: result.stderr || '',
                programOutput: result.stdout || '' // AST output goes to stdout
            })
        } else {

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
    console.log('📋 RUN REQUEST:')
    console.log(code || '(no code)')

    if (!code) {
        return res.status(400).json({ success: false, error: 'No code provided' })
    }

    try {
        const result = await runOspreyCompiler(['--sandbox', '--run'], code)

        if (result.success) {
            res.status(200).json({
                success: true,
                compilerOutput: result.stderr || '',
                programOutput: result.stdout || ''
            })
        } else {
            const errorOutput = result.stderr || result.stdout || '';

            // Detect INTERNAL compiler errors
            const isInternalError = errorOutput.includes('INTERNAL_COMPILER_ERROR:');
            if (isInternalError) {
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

            const statusCode = isCompilationError ? 422 : 400;

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
        await fs.rm(tempBaseDir, { recursive: true, force: true })
    } catch (error) {
        // Ignore errors
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
            }
        }
    } catch (error) {
        // Ignore errors
    }
}

// Run cleanup every 30 minutes
setInterval(cleanupOldTempFolders, 30 * 60 * 1000)

// DELETE ALL TEMP FOLDERS ON STARTUP
deleteAllTempFolders()

// Helper function to run Osprey compiler
// Always uses --sandbox flag for security (disables HTTP, WebSocket, file system, and FFI access)
function runOspreyCompiler(args, code = '') {
    return new Promise(async (resolve, reject) => {
        // Simple temp file approach - use absolute path for Docker
        const tempBaseDir = process.env.DOCKER_ENV
            ? '/app/temp'
            : path.resolve(process.cwd(), 'temp')
        const tempFile = path.join(tempBaseDir, 'main.osp')

        console.log(`🔧 Setting up temp file: ${tempFile}`)

        try {
            // Make sure the temp directory exists
            await fs.mkdir(tempBaseDir, { recursive: true, mode: 0o755 })

            // Write the code file
            await fs.writeFile(tempFile, code, { mode: 0o644 })

            // Use the osprey binary from PATH (installed in Docker) or fallback to local dev path
            const ospreyPath = process.env.NODE_ENV === 'production' || process.env.DOCKER_ENV
                ? 'osprey'
                : path.resolve(__dirname, '../../compiler/bin/osprey')

            // Set working directory so compiler can find runtime libraries
            const workingDir = process.env.NODE_ENV === 'production' || process.env.DOCKER_ENV
                ? process.cwd() // Docker: osprey binary in /usr/local/bin uses FHS ../lib = /usr/local/lib
                : path.resolve(__dirname, '../../compiler') // Dev: osprey binary in compiler/bin uses ../lib = compiler/lib

            console.log(`🔧 Running: ${ospreyPath} ${tempFile} ${args.join(' ')}`)

            const child = spawn(ospreyPath, [tempFile, ...args], {
                stdio: 'pipe',
                cwd: workingDir,
                timeout: 5000,
                env: {
                    ...process.env,
                    TMPDIR: tempBaseDir,
                    TMP: tempBaseDir,
                    TEMP: tempBaseDir
                }
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
                console.log('📤 COMPILER OUTPUT:')
                console.log(stderr || '(no stderr)')
                console.log('📤 PROGRAM OUTPUT:')
                console.log(stdout || '(no stdout)')

                // Clean up the temp file
                try {
                    await fs.rm(tempFile, { force: true })
                    console.log(`🧹 Cleaned up: ${tempFile}`)
                } catch (e) {
                    console.warn('⚠️ Failed to cleanup:', e.message)
                }

                resolve({
                    exitCode,
                    stdout,
                    stderr,
                    success: exitCode === 0
                })
            })

            child.on('error', async (error) => {
                console.error('❌ Spawn error:', error.message)
                // Clean up temp file on error
                try {
                    await fs.rm(tempFile, { force: true })
                } catch (e) {
                    // Ignore cleanup errors
                }
                reject(error)
            })
        } catch (error) {
            console.error('❌ Setup error:', error.message)
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