---
layout: base.njk
title: "Osprey Playground"
description: "Try Osprey programming language online with interactive code examples and real-time compilation"
---

<link rel="stylesheet" data-name="vs/editor/editor.main" href="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs/editor/editor.main.min.css">

<style>
    /* Override website layout constraints for playground area */
    .main-content {
        padding: 0 !important;
        margin: 0 !important;
        max-width: none !important;
    }
    
    .playground-container {
        display: flex;
        flex-direction: column;
        background: #1e1e1e;
        color: #d4d4d4;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        min-height: calc(100vh - 80px);
        height: calc(100vh - 80px);
    }
    
    .main {
        display: flex;
        flex: 1;
        overflow: hidden;
        min-height: 0;
    }
    
    .editor-container {
        flex: 1;
        display: flex;
        flex-direction: column;
        min-height: 0;
    }
    
    .editor-header {
        background: #2d2d30;
        padding: 10px 20px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-bottom: 1px solid #444;
        flex-shrink: 0;
    }
    
    .editor-title {
        display: flex;
        align-items: center;
        gap: 10px;
        font-size: 14px;
    }
    
    .playground-badge {
        font-size: 12px;
        color: #569cd6;
        opacity: 0.8;
    }
    
    .header-right {
        display: flex;
        align-items: center;
        gap: 15px;
    }
    
    .status {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 12px;
    }
    
    .status-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #ffa500;
    }
    
    .status-dot.connected {
        background: #4ec9b0;
    }
    
    .status-dot.error {
        background: #f44747;
    }
    
    .button-group {
        display: flex;
        gap: 0;
    }
    
    #editor {
        flex: 1;
        min-height: 0;
        height: 100%;
    }
    
    .output-container {
        width: 400px;
        display: flex;
        flex-direction: column;
        border-left: 1px solid #444;
        min-height: 0;
    }
    
    .output-header {
        background: #2d2d30;
        padding: 10px 20px;
        border-bottom: 1px solid #444;
        display: flex;
        justify-content: space-between;
        align-items: center;
        flex-shrink: 0;
    }
    
    #output {
        flex: 1;
        padding: 20px;
        overflow-y: auto;
        font-family: 'Consolas', 'Monaco', monospace;
        white-space: pre-wrap;
        min-height: 0;
    }
    
    #output.error {
        color: #f44747;
    }
    
    /* Splitter styles */
    .splitter {
        background: #444;
        cursor: col-resize;
        position: relative;
        flex-shrink: 0;
        width: 4px;
        transition: background-color 0.2s ease;
    }
    
    .splitter:hover {
        background: #569cd6;
    }
    
    .splitter::before {
        content: '';
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 2px;
        height: 20px;
        background: #666;
        border-radius: 1px;
    }
    
    .splitter.dragging {
        background: #569cd6;
    }
    
    /* Mobile responsiveness */
    @media (max-width: 768px) {
        .playground-container {
            height: 100vh;
            min-height: 100vh;
        }
        
        .main {
            flex-direction: column;
        }
        
        .editor-container {
            flex: 1;
        }
        
        .output-container {
            width: 100%;
            height: 40%;
            border-left: none;
            border-top: 1px solid #444;
        }
        
        .splitter {
            cursor: row-resize;
            width: 100%;
            height: 4px;
            border-top: none;
        }
        
        .splitter::before {
            width: 20px;
            height: 2px;
        }
        
        .editor-header {
            padding: 8px 15px;
        }
        
        .header-right {
            gap: 10px;
        }
        
        .editor-title {
            gap: 5px;
            font-size: 13px;
        }
        
        .playground-badge {
            display: none;
        }
        
        .status {
            gap: 5px;
            font-size: 11px;
        }
        
        button {
            padding: 6px 12px;
            font-size: 13px;
            margin-left: 5px;
        }
        
        .output-header {
            padding: 8px 15px;
        }
        
        #output {
            padding: 15px;
        }
    }
    
    @media (max-width: 480px) {
        .editor-header, .output-header {
            padding: 6px 10px;
        }
        
        .header-right {
            gap: 8px;
        }
        
        .editor-title {
            font-size: 12px;
        }
        
        .status {
            font-size: 10px;
        }
        
        button {
            padding: 5px 8px;
            font-size: 12px;
            margin-left: 3px;
        }
        
        #output {
            padding: 10px;
            font-size: 13px;
        }
        

    }
    
    button {
        background: #0e639c;
        color: white;
        border: none;
        padding: 8px 16px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 14px;
        margin-left: 10px;
    }
    
    button:hover {
        background: #1177bb;
    }
    
    button.primary {
        background: #16825d;
    }
    
    button.primary:hover {
        background: #1ea571;
    }
</style>

<div class="playground-container">
    <div class="main">
        <div class="editor-container">
            <div class="editor-header">
                <div class="editor-title">
                    <span>Osprey Editor</span>
                    <span class="playground-badge">⚡ Playground</span>
                </div>
                <div class="header-right">
                    <div class="status">
                        <div id="status-dot" class="status-dot"></div>
                        <span id="status-text">Connecting...</span>
                    </div>
                    <div class="button-group">
                        <button onclick="compileCode()">Compile</button>
                        <button class="primary" onclick="runCode()">Run</button>
                    </div>
                </div>
            </div>
            <div id="editor"></div>
        </div>
        
        <div id="splitter" class="splitter"></div>
        
        <div class="output-container">
            <div class="output-header">
                <span>Output</span>
                <button onclick="clearOutput()">Clear</button>
            </div>
            <div id="output"></div>
        </div>
    </div>
</div>

<!-- Load Monaco from CDN -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs/loader.min.js"></script>

<script>
    let editor;
    const API_URL = 'https://osprey-web-compiler-gateway.mail-bff.workers.dev/api';
    
    // Initialize Monaco Editor
    require.config({ paths: { vs: 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs' } });
    
    require(['vs/editor/editor.main'], function() {
        // Register Osprey language
        monaco.languages.register({ id: 'osprey' });
        
        // Define syntax highlighting
        monaco.languages.setMonarchTokensProvider('osprey', {
            keywords: ['fn', 'let', 'mut', 'type', 'import', 'match', 'if', 'else', 'loop', 'spawn', 'extern', 'true', 'false'],
            tokenizer: {
                root: [
                    [/\/\/.*$/, 'comment'],
                    [/[a-z_$][\w$]*/, {
                        cases: {
                            '@keywords': 'keyword',
                            '@default': 'identifier'
                        }
                    }],
                    [/".*?"/, 'string'],
                    [/\d+/, 'number'],
                ]
            }
        });
        
        // Create editor
        editor = monaco.editor.create(document.getElementById('editor'), {
            value: `// Simple Osprey Demo - Basic constructs that definitely work
// Pattern matching, functional pipes, and string interpolation

fn double(x: int) -> int = x * 2
fn add10(x: int) -> int = x + 10

fn gradeScore(score: int) -> string = match score {
    100 => "Perfect"
    90 => "Excellent" 
    80 => "Good"
    70 => "Fair"
    _ => "Needs work"
}

print("=== Osprey Language Demo ===")

// Basic functional pipeline
let result = 5 |> double |> add10
print("Pipeline: 5 -> double -> add10 = \${result}")

// Pattern matching demo
let score1 = 100
let score2 = 85
let score3 = 60

print("Scores:")
print("\${score1}: \${gradeScore(score1)}")
print("\${score2}: \${gradeScore(score2)}")
print("\${score3}: \${gradeScore(score3)}")

// Math operations
let a = 15
let b = 25
let sum = a + b
let product = a * b

print("Math: \${a} + \${b} = \${sum}")
print("Math: \${a} * \${b} = \${product}")

// Range iteration
print("Numbers 1-5:")
range(1, 6) |> forEach(print)

print("Demo complete!")`,
            language: 'osprey',
            theme: 'vs-dark',
            automaticLayout: true
        });
        
        // Update status
        updateStatus('connected', 'Ready');
    });
    
    function updateStatus(type, message) {
        const statusDot = document.getElementById('status-dot');
        const statusText = document.getElementById('status-text');
        
        statusDot.className = `status-dot ${type}`;
        statusText.textContent = message;
    }

    async function compileCode() {
        const code = editor.getValue();
        const output = document.getElementById('output');
        
        updateStatus('', 'Compiling...');
        output.textContent = 'Compiling...';
        
        try {
            const response = await fetch(`${API_URL}/compile`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ code })
            });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            const result = await response.json();
            
            output.className = result.success ? '' : 'error';
            output.textContent = result.success ? result.output : result.error;
            updateStatus('connected', 'Ready');
            
        } catch (error) {
            output.className = 'error';
            output.textContent = `Compilation failed: ${error.message}`;
            updateStatus('error', 'Error');
        }
    }
    
    async function runCode() {
        const code = editor.getValue();
        const output = document.getElementById('output');
        
        updateStatus('', 'Running...');
        output.textContent = 'Running...';
        
        try {
            const response = await fetch(`${API_URL}/run`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ code })
            });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            const result = await response.json();
            
            output.className = result.success ? '' : 'error';
            output.textContent = result.success ? result.output : result.error;
            updateStatus('connected', 'Ready');
            
        } catch (error) {
            output.className = 'error';
            output.textContent = `Execution failed: ${error.message}`;
            updateStatus('error', 'Error');
        }
    }
    
    function clearOutput() {
        document.getElementById('output').textContent = '';
        document.getElementById('output').className = '';
    }
    
    // Splitter functionality
    let isDragging = false;
    let startX = 0;
    let startY = 0;
    let startWidth = 0;
    let startHeight = 0;
    let isMobile = false;
    
    function initSplitter() {
        const splitter = document.getElementById('splitter');
        const editorContainer = document.querySelector('.editor-container');
        const outputContainer = document.querySelector('.output-container');
        
        if (!splitter || !editorContainer || !outputContainer) return;
        
        splitter.addEventListener('mousedown', startDrag);
        document.addEventListener('mousemove', drag);
        document.addEventListener('mouseup', stopDrag);
        
        // Touch events for mobile
        splitter.addEventListener('touchstart', startDrag);
        document.addEventListener('touchmove', drag);
        document.addEventListener('touchend', stopDrag);
        
        // Check if mobile layout
        function checkMobile() {
            isMobile = window.innerWidth <= 768;
        }
        
        checkMobile();
        window.addEventListener('resize', checkMobile);
    }
    
    function startDrag(e) {
        isDragging = true;
        const splitter = document.getElementById('splitter');
        const editorContainer = document.querySelector('.editor-container');
        const outputContainer = document.querySelector('.output-container');
        
        splitter.classList.add('dragging');
        
        if (isMobile) {
            startY = e.touches ? e.touches[0].clientY : e.clientY;
            startHeight = editorContainer.offsetHeight;
        } else {
            startX = e.touches ? e.touches[0].clientX : e.clientX;
            startWidth = editorContainer.offsetWidth;
        }
        
        e.preventDefault();
    }
    
    function drag(e) {
        if (!isDragging) return;
        
        const main = document.querySelector('.main');
        const editorContainer = document.querySelector('.editor-container');
        const outputContainer = document.querySelector('.output-container');
        
                 if (isMobile) {
             const currentY = e.touches ? e.touches[0].clientY : e.clientY;
             const deltaY = currentY - startY;
             const newHeight = startHeight + deltaY;
             const mainHeight = main.offsetHeight;
             
             if (newHeight >= 0 && newHeight <= mainHeight) {
                 const heightPercent = (newHeight / mainHeight) * 100;
                 const outputPercent = 100 - heightPercent;
                 
                 editorContainer.style.flex = 'none';
                 editorContainer.style.height = `${heightPercent}%`;
                 outputContainer.style.height = `${outputPercent}%`;
             }
         } else {
             const currentX = e.touches ? e.touches[0].clientX : e.clientX;
             const deltaX = currentX - startX;
             const newWidth = startWidth + deltaX;
             const mainWidth = main.offsetWidth;
             
             if (newWidth >= 0 && newWidth <= mainWidth) {
                 const widthPercent = (newWidth / mainWidth) * 100;
                 const outputWidth = mainWidth - newWidth - 4; // Account for splitter width
                 
                 editorContainer.style.flex = 'none';
                 editorContainer.style.width = `${newWidth}px`;
                 outputContainer.style.width = `${outputWidth}px`;
             }
         }
        
        e.preventDefault();
    }
    
    function stopDrag() {
        if (!isDragging) return;
        
        isDragging = false;
        const splitter = document.getElementById('splitter');
        splitter.classList.remove('dragging');
    }
    
    // Initialize splitter when page loads
    window.addEventListener('load', initSplitter);
</script> 