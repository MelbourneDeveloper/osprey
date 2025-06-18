import { spawn } from 'child_process'
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
        version: '1.0.0',
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
        console.log('✅ Compile success, output length:', result.stdout.length)
        res.json({ success: true, output: result.stdout })
    } catch (error) {
        console.error('❌ Compile error:', error.message)
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
        const result = await runOspreyCompiler(['--sandbox', '--run'], code)
        console.log('✅ Run success, output length:', result.stdout.length)
        res.json({ success: true, output: result.stdout })
    } catch (error) {
        console.error('❌ Run error:', error.message)
        res.status(500).json({ success: false, error: error.message })
    }
})

// Helper function to run Osprey compiler
// Always uses --sandbox flag for security (disables HTTP, WebSocket, file system, and FFI access)
function runOspreyCompiler(args, code = '') {
    return new Promise(async (resolve, reject) => {
        // Ensure the temp directory exists
        const tempDir = '/tmp/osprey-temp'
        const tempFile = path.join(tempDir, `temp_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.osp`)

        try {
            // Create the temp directory if it doesn't exist
            await fs.mkdir(tempDir, { recursive: true })

            console.log(`💾 Writing temp file: ${tempFile}`)
            await fs.writeFile(tempFile, code)

            // Use the osprey binary from PATH (installed in Docker) or fallback to local dev path
            const ospreyPath = process.env.NODE_ENV === 'production' || process.env.DOCKER_ENV
                ? 'osprey'
                : path.resolve(__dirname, '../../compiler/bin/osprey')
            console.log(`🔨 Running: ${ospreyPath} ${tempFile} ${args.join(' ')}`)
            const child = spawn(ospreyPath, [tempFile, ...args], {
                stdio: 'pipe',
                cwd: process.cwd(),
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

            child.on('close', async (code) => {
                // Clean up temp file
                try {
                    await fs.unlink(tempFile)
                    console.log(`🗑️ Cleaned up temp file: ${tempFile}`)
                } catch (e) {
                    console.error('⚠️ Failed to clean up temp file:', e.message)
                }

                if (code === 0) {
                    resolve({ stdout, stderr })
                } else {
                    reject(new Error(stderr || stdout || `Process exited with code ${code}`))
                }
            })

            child.on('error', async (error) => {
                // Clean up temp file on error
                try {
                    await fs.unlink(tempFile)
                } catch (e) {
                    // Ignore cleanup errors
                }
                reject(error)
            })
        } catch (error) {
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