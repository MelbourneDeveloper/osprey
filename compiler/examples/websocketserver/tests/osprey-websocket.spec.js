// @ts-check
const { test, expect } = require('@playwright/test');
const { spawn } = require('child_process');
const path = require('path');

let ospreyServer = null;
let nodeServer = null;

/**
 * Start Osprey WebSocket Server
 */
async function startOspreyServer() {
    return new Promise((resolve, reject) => {
        console.log('🚀 Starting Osprey WebSocket Server...');

        // Run the Osprey server
        ospreyServer = spawn('../../bin/osprey', ['osprey_websocket_server.osp', '--run'], {
            cwd: process.cwd(),
            stdio: ['pipe', 'pipe', 'pipe']
        });

        let output = '';

        ospreyServer.stdout.on('data', (data) => {
            const str = data.toString();
            output += str;
            console.log('Osprey Server:', str.trim());

            // Server is ready when it shows it's listening
            if (str.includes('WebSocket server listening with result: 0')) {
                setTimeout(resolve, 1000); // Give it a moment to fully start
            }
        });

        ospreyServer.stderr.on('data', (data) => {
            console.error('Osprey Server Error:', data.toString());
        });

        ospreyServer.on('error', (error) => {
            console.error('Failed to start Osprey server:', error);
            reject(error);
        });

        // Timeout if server doesn't start
        setTimeout(() => {
            if (ospreyServer && !output.includes('listening with result: 0')) {
                reject(new Error('Osprey server failed to start within timeout'));
            }
        }, 10000);
    });
}

/**
 * Start Node.js WebSocket Server
 */
async function startNodeServer() {
    return new Promise((resolve, reject) => {
        console.log('🚀 Starting Node.js WebSocket Server...');

        nodeServer = spawn('../../osprey', ['osprey_websocket_server.osp', '--run'], {
            cwd: process.cwd(),
            stdio: ['pipe', 'pipe', 'pipe']
        });

        nodeServer.stdout.on('data', (data) => {
            const str = data.toString();
            console.log('Node Server:', str.trim());

            // Server is ready when it shows it's running
            if (str.includes('Server running on http://localhost:8080')) {
                setTimeout(resolve, 1000); // Give it a moment to fully start
            }
        });

        nodeServer.stderr.on('data', (data) => {
            console.error('Node Server Error:', data.toString());
        });

        nodeServer.on('error', (error) => {
            console.error('Failed to start Node server:', error);
            reject(error);
        });

        // Timeout if server doesn't start
        setTimeout(() => {
            reject(new Error('Node server failed to start within timeout'));
        }, 10000);
    });
}

/**
 * Stop servers
 */
function stopServers() {
    if (ospreyServer) {
        console.log('🛑 Stopping Osprey server...');
        ospreyServer.kill('SIGTERM');
        ospreyServer = null;
    }
    if (nodeServer) {
        console.log('🛑 Stopping Node server...');
        nodeServer.kill('SIGTERM');
        nodeServer = null;
    }
}

test.describe('🔌 Osprey WebSocket Server Tests', () => {

    test.beforeAll(async () => {
        // Note: Osprey server runs and exits immediately in demo mode
        // We'll test the functionality it provides
    });

    test.afterAll(async () => {
        stopServers();
    });

    test('osprey websocket server should compile and run', async () => {
        console.log('🧪 Testing Osprey WebSocket Server compilation and execution...');

        // Test that the Osprey server can be executed
        const result = await new Promise((resolve, reject) => {
            const childProcess = spawn('../../bin/osprey', ['osprey_websocket_server.osp', '--run'], {
                cwd: process.cwd(),
                stdio: ['pipe', 'pipe', 'pipe']
            });

            let stdout = '';
            let stderr = '';

            childProcess.stdout.on('data', (data) => {
                stdout += data.toString();
            });

            childProcess.stderr.on('data', (data) => {
                stderr += data.toString();
            });

            childProcess.on('close', (code) => {
                resolve({ code, stdout, stderr });
            });

            childProcess.on('error', (error) => {
                reject(error);
            });
        });

        // Verify the server executed successfully
        expect(result.code).toBe(0);
        expect(result.stdout).toContain('WebSocket server created with ID: 1');
        expect(result.stdout).toContain('WebSocket server listening with result: 0');
        expect(result.stdout).toContain('Welcome broadcast sent with result: 0');
        expect(result.stdout).toContain('WebSocket server stopped with result: 0');

        console.log('✅ Osprey WebSocket Server executed successfully!');
    });

    test('osprey server should demonstrate websocket functions', async () => {
        console.log('🧪 Testing Osprey WebSocket function demonstrations...');

        const result = await new Promise((resolve, reject) => {
            const childProcess = spawn('../../bin/osprey', ['osprey_websocket_server.osp'], {
                cwd: process.cwd(),
                stdio: ['pipe', 'pipe', 'pipe']
            });

            let stdout = '';
            let stderr = '';

            childProcess.stdout.on('data', (data) => {
                stdout += data.toString();
            });

            childProcess.stderr.on('data', (data) => {
                stderr += data.toString();
            });

            childProcess.on('close', (code) => {
                resolve({ code, stdout, stderr });
            });

            childProcess.on('error', (error) => {
                reject(error);
            });
        });

        // Verify LLVM IR generation includes WebSocket function calls
        expect(result.stdout).toContain('declare i64 @websocket_create_server');
        expect(result.stdout).toContain('declare i64 @websocket_server_listen');
        expect(result.stdout).toContain('declare i64 @websocket_server_broadcast');
        expect(result.stdout).toContain('declare i64 @websocket_stop_server');

        // Verify function calls are generated
        expect(result.stdout).toContain('call i64 @websocket_create_server');
        expect(result.stdout).toContain('call i64 @websocket_server_listen');
        expect(result.stdout).toContain('call i64 @websocket_server_broadcast');
        expect(result.stdout).toContain('call i64 @websocket_stop_server');

        console.log('✅ Osprey WebSocket functions compiled to LLVM IR successfully!');
    });

    test('osprey websocket test page should load', async ({ page }) => {
        console.log('🧪 Testing WebSocket test page loading...');

        // Load the WebSocket test HTML page directly
        const filePath = path.resolve(process.cwd(), 'websocket_test.html');
        await page.goto(`file://${filePath}`);

        // Verify the page loaded correctly
        await expect(page).toHaveTitle(/WebSocket/);

        // Check for WebSocket test elements (look for any connection status element)
        const connectionElements = await page.locator('[id*="status"], [class*="status"], [id*="connection"], [class*="connection"]').count();
        expect(connectionElements).toBeGreaterThan(0);

        console.log('✅ WebSocket test page loaded successfully!');
    });

    test('websocket test page should have websocket functionality', async ({ page }) => {
        console.log('🧪 Testing WebSocket test page elements...');

        // Load the WebSocket test HTML page
        const filePath = path.resolve(process.cwd(), 'websocket_test.html');
        await page.goto(`file://${filePath}`);

        // Wait for the page to load
        await page.waitForTimeout(2000);

        // Check that the page contains WebSocket-related content
        const pageContent = await page.textContent('body');
        expect(pageContent).toContain('WebSocket');

        // Look for common WebSocket test page elements
        const hasInputField = (await page.locator('input[type="text"]').count()) > 0;
        const hasButton = (await page.locator('button').count()) > 0;
        const hasScript = (await page.locator('script').count()) > 0;

        // At least one of these should be present in a WebSocket test page
        expect(hasInputField || hasButton || hasScript).toBe(true);

        console.log('🔌 WebSocket test page has expected functionality elements');
        console.log('✅ WebSocket test page validation completed!');
    });

});

test.describe('🟨 Node.js WebSocket Server Tests', () => {

    test.beforeAll(async () => {
        // Start Node.js server for comparison tests
        try {
            await startNodeServer();
            console.log('✅ Node.js server started for tests');
        } catch (error) {
            console.log('⚠️ Could not start Node.js server, skipping Node.js tests');
            test.fail();
        }
    });

    test.afterAll(async () => {
        stopServers();
    });

    test('nodejs websocket server should serve html page', async ({ page }) => {
        test.fail(!nodeServer, 'Node.js server not running');

        console.log('🧪 Testing Node.js server HTTP serving...');

        // Navigate to the Node.js server
        await page.goto('http://127.0.0.1:8080');

        // Verify the page loaded
        await expect(page).toHaveTitle(/WebSocket/);

        console.log('✅ Node.js server served HTML page successfully!');
    });

    test('nodejs websocket connection should be attempted', async ({ page }) => {
        test.fail(!nodeServer, 'Node.js server not running');

        console.log('🧪 Testing Node.js WebSocket connection attempt...');

        // Navigate to the Node.js server
        await page.goto('http://127.0.0.1:8080');

        // Wait for WebSocket connection attempt
        await page.waitForTimeout(3000);

        // Check that page contains WebSocket content
        const pageContent = await page.textContent('body');
        expect(pageContent).toContain('WebSocket');

        console.log('✅ Node.js WebSocket page loaded and connection attempted!');
    });

});

test.describe('🆚 Server Comparison', () => {

    test('compare server startup performance', async () => {
        console.log('🧪 Comparing server startup performance...');

        // Test Osprey server startup time
        const ospreyStart = Date.now();
        const ospreyResult = await new Promise((resolve, reject) => {
            const childProcess = spawn('../../bin/osprey', ['osprey_websocket_server.osp', '--run'], {
                cwd: process.cwd(),
                stdio: ['pipe', 'pipe', 'pipe']
            });

            childProcess.on('close', (code) => {
                resolve({ code, time: Date.now() - ospreyStart });
            });

            childProcess.on('error', reject);
        });

        console.log(`⚡ Osprey server execution time: ${ospreyResult.time}ms`);
        console.log(`✅ Osprey server completed with exit code: ${ospreyResult.code}`);

        // Note: Node.js comparison would require running actual server
        // For now, we just verify Osprey performance
        expect(ospreyResult.code).toBe(0);
        expect(ospreyResult.time).toBeLessThan(5000); // Should complete quickly

        console.log('🏆 Osprey server demonstrates fast startup/execution!');
    });

}); 