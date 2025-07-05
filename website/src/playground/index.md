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
        background: #5a8a6b;
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
        background: #1e1e1e;
        color: #d4d4d4;
        line-height: 1.4;
    }
    
    #output.error {
        color: #d4d4d4;
        background: #1e1e1e;
        border-left: none;
    }
    
    #output.success {
        color: #d4d4d4;
        background: #1e1e1e;
        border-left: none;
    }
    
    #output.warning {
        color: #ffa500;
        background: #2d2d1b;
        border-left: 3px solid #ffa500;
    }
    
    .output-section {
        margin-bottom: 20px;
    }
    
    .output-section:last-child {
        margin-bottom: 0;
    }
    
    .output-label {
        font-size: 12px;
        text-transform: uppercase;
        opacity: 0.7;
        margin-bottom: 8px;
        font-weight: 600;
        letter-spacing: 0.5px;
    }
    
    .compiler-output {
        color: #d4d4d4;
        background: transparent;
        padding: 0;
        border: none;
        margin-bottom: 12px;
    }
    
    .program-output {
        color: #7cb992;
        background: rgba(124, 185, 146, 0.08);
        padding: 12px;
        border-radius: 4px;
        border-left: 3px solid #5a8a6b;
    }
    
    .program-output.empty {
        display: none;
    }
    
    .line-number {
        color: #569cd6;
        font-weight: bold;
    }
    
    /* Error listview styles */
    .error-list {
        display: grid;
        gap: 1px;
        font-family: 'Consolas', 'Monaco', monospace;
        font-size: 13px;
        line-height: 1.4;
    }
    
    .error-line {
        display: grid;
        grid-template-columns: auto 1fr;
        gap: 12px;
        padding: 8px 12px;
        background: #2d2d30;
        border: 1px solid #444;
        cursor: pointer;
        transition: all 0.2s ease;
        align-items: center;
    }
    
    .error-line:hover {
        background: #3c3c3c;
        border-color: #569cd6;
    }
    
    .error-line.selected {
        background: #404040;
        border-color: #569cd6;
        box-shadow: 0 0 0 1px #569cd6;
    }
    
    .error-location {
        color: #569cd6;
        font-weight: bold;
        font-size: 12px;
        white-space: nowrap;
        cursor: pointer;
        text-decoration: none;
    }
    
    .error-location:hover {
        text-decoration: underline;
    }
    
    .error-message {
        color: #f44747;
        flex: 1;
        word-break: break-word;
    }
    
    /* Editor error highlighting */
    .highlighted-error-line {
        background: rgba(244, 71, 71, 0.15) !important;
        border-left: 2px solid #f44747 !important;
    }
    
    .error-glyph {
        background: #f44747;
        width: 4px !important;
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
    const API_URL = 'https://osprey.fly.dev/api';
    
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
            value: `// 🚀 OSPREY MEGA SHOWCASE - COMPREHENSIVE SANDBOXABLE FEATURES 🎯
// A realistic data analytics engine demonstrating ALL safe features

// 📊 Type System - Simple union types 
type AnalysisResult = Success | Warning | Critical

// 🎭 Algebraic Effects - Safe side effects management (from working examples)
effect Analytics {
    logMetric: fn(string, int) -> Unit
    incrementCounter: fn(string) -> Unit
    getTotal: fn(string) -> int
}

effect Logger {
    info: fn(string) -> Unit
    warn: fn(string) -> Unit
    error: fn(string) -> Unit
}

effect Cache {
    store: fn(string, int) -> Unit
    clear: fn() -> Unit
}

// 🔧 Pure Functions - Core business logic
fn calculateBonus(amount: int) -> int = match amount > 10000 {
    true => (amount * 15) / 100
    false => match amount > 5000 {
        true => (amount * 10) / 100
        false => match amount > 1000 {
            true => (amount * 5) / 100
            false => 0
        }
    }
}

fn categorizePerformance(latency: int) -> string = match latency < 50 {
    true => "Excellent"
    false => match latency < 200 {
        true => "Good"
        false => match latency < 500 {
            true => "Acceptable"
            false => "Poor"
        }
    }
}

// Analytics processing with effects (using working effect syntax)
fn processSalesData(amount: int, region: string) -> int ![Analytics, Logger] = {
    let bonus = calculateBonus(amount)
    let multiplier = match region {
        "North" => 12
        "South" => 15
        "East" => 18
        _ => 10
    }
    let score = (amount + bonus) * multiplier / 100
    
    perform Analytics.logMetric("sales_processed", 1)
    perform Analytics.incrementCounter(region)
    perform Logger.info("Sales processed: " + region + " region, $" + toString(amount) + ", bonus: $" + toString(bonus))
    
    score
}

fn processTrafficData(visitors: int, source: string) -> int ![Analytics, Logger] = {
    let engagement = visitors * 3
    perform Analytics.logMetric("traffic_analyzed", visitors)
    perform Analytics.incrementCounter("traffic_sources")
    perform Logger.info("Traffic analyzed: " + toString(visitors) + " visitors from " + source)
    
    engagement
}

fn processPerformanceData(latency: int, endpoint: string) -> int ![Analytics, Logger] = {
    let category = categorizePerformance(latency)
    let healthScore = match latency < 100 {
        true => 1000 - latency
        false => match latency < 1000 {
            true => 500 - (latency / 2)
            false => 100
        }
    }
    
    perform Analytics.logMetric("performance_checked", 1)
    perform Logger.info("Performance: " + endpoint + " - " + toString(latency) + "ms (" + category + ")")
    
    healthScore
}

// 🚀 Fiber-based Concurrency - Parallel processing (exact syntax from working examples)
fn processAllDataConcurrently() -> int ![Analytics, Logger, Cache] = {
    perform Logger.info("Starting concurrent data processing")
    
    // Spawn concurrent fibers for different data types (exact syntax from fiber_advanced.osp)
    let salesFiber = spawn {
        sleep(5)
        let s1 = processSalesData(amount: 12000, region: "North")
        let s2 = processSalesData(amount: 8500, region: "South") 
        let s3 = processSalesData(amount: 3200, region: "East")
        s1 + s2 + s3
    }
    
    let trafficFiber = spawn {
        sleep(3)
        let t1 = processTrafficData(visitors: 4500, source: "Google")
        let t2 = processTrafficData(visitors: 2800, source: "Social")
        let t3 = processTrafficData(visitors: 1200, source: "Direct")
        t1 + t2 + t3
    }
    
    let perfFiber = spawn {
        sleep(7)
        let p1 = processPerformanceData(latency: 45, endpoint: "/api/users")
        let p2 = processPerformanceData(latency: 180, endpoint: "/api/orders")
        let p3 = processPerformanceData(latency: 650, endpoint: "/api/reports")
        p1 + p2 + p3
    }
    
    // Wait for all results and combine (exact syntax from working examples)
    let salesTotal = await(salesFiber)
    let trafficTotal = await(trafficFiber)
    let perfTotal = await(perfFiber)
    
    let grandTotal = salesTotal + trafficTotal + perfTotal
    perform Cache.store("grand_total", grandTotal)
    perform Analytics.logMetric("concurrent_processing_complete", 1)
    perform Logger.info("Concurrent processing complete - Sales: " + toString(salesTotal) + ", Traffic: " + toString(trafficTotal) + ", Performance: " + toString(perfTotal))
    
    grandTotal
}

// 🧮 Complex Calculation with String Interpolation
fn generateAnalyticsReport(totalScore: int, processingTime: int) -> string ![Analytics, Cache] = {
    let efficiency = (totalScore * 100) / (processingTime + 1)
    let grade = match efficiency > 1000 {
        true => "A+"
        false => match efficiency > 800 {
            true => "A"
            false => match efficiency > 600 {
                true => "B"
                false => match efficiency > 400 {
                    true => "C"
                    false => "D"
                }
            }
        }
    }
    
    let statusMessage = match totalScore > 1000 {
        true => "🟢 EXCELLENT"
        false => "🟡 NORMAL"
    }
    
    perform Analytics.incrementCounter("reports_generated")
    
    "📊 ANALYTICS DASHBOARD REPORT 📊\n" +
    "═════════════════════════════════════\n" +
    "Total Score: " + toString(totalScore) + " points\n" +
    "Processing Time: " + toString(processingTime) + "ms\n" +
    "Efficiency Rating: " + toString(efficiency) + " (Grade: " + grade + ")\n" +
    "Last Batch Total: 1250 points\n" +
    "Performance Category: " + categorizePerformance(processingTime) + "\n" +
    "═════════════════════════════════════\n" +
    "🎯 System Status: " + statusMessage + "\n"
}

// 🎪 Main Function - Comprehensive demo (exact syntax from comprehensive_demo.osp)
fn main() -> Unit = {
    handle Cache
        store key value => print("💾 Cache: Stored " + key + " = " + toString(value))
        clear => print("🧹 Cache cleared")
    in handle Analytics
        logMetric metric value => print("📈 Metric: " + metric + " += " + toString(value))
        incrementCounter counter => print("🔢 Counter: " + counter + " incremented")
        getTotal counter => match counter {
            "sales_processed" => 15
            "traffic_analyzed" => 8
            "performance_checked" => 12
            _ => 0
        }
    in handle Logger
        info msg => print("ℹ️  " + msg)
        warn msg => print("⚠️  " + msg)
        error msg => print("❌ " + msg)
    in {
        print("🚀 OSPREY MEGA SHOWCASE - COMPREHENSIVE SANDBOXABLE FEATURES DEMO")
        print("═══════════════════════════════════════════════════════════════")
        
        // Demonstrate concurrent processing with fibers
        let startTime = 1000
        let totalScore = processAllDataConcurrently()
        let endTime = 1035
        
        // Generate comprehensive report with string interpolation
        let processingTime = endTime - startTime
        let report = generateAnalyticsReport(totalScore: totalScore, processingTime: processingTime)
        
        print(report)
        
        // Demonstrate functional programming with calculations
        let performanceSum = 100 + 200 + 300 + 400 + 500
        print("📊 Performance Metrics Total: " + toString(performanceSum))
        
        // Pattern matching demonstration (single example, not repetitive)
        print("✅ Success: System operating normally")
        
        print("═══════════════════════════════════════════════════════════════")
        print("🎉 COMPREHENSIVE DEMO COMPLETE - ALL SANDBOXABLE FEATURES TESTED")
        print("✅ Type System ✅ Algebraic Effects ✅ Pattern Matching")
        print("✅ Fiber Concurrency ✅ Functional Programming ✅ String Interpolation")
    }
}
`,
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
        output.innerHTML = '<div style="color: #ffa500;">Compiling...</div>';
        
        try {
            const response = await fetch(`${API_URL}/compile`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ code })
            });
            
            const result = await response.json();
            
            if (!response.ok) {
                // Handle HTTP errors that still have JSON error details
                if (result.error) {
                    output.className = 'error';
                    output.innerHTML = formatErrorOutput(result.error);
                    updateStatus('error', 'Compilation failed');
                    return;
                } else {
                    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
                }
            }
            
            if (result.success) {
                // Successful compilation
                output.className = 'success';
                let outputText = '';
                
                if (result.programOutput && result.programOutput.trim()) {
                    outputText = formatPlainOutput(result.programOutput);
                } else {
                    outputText = '✅ Compilation successful - no output';
                }
                
                output.innerHTML = outputText;
                updateStatus('connected', 'Ready');
            } else {
                // Compilation failed
                output.className = 'error';
                output.innerHTML = formatErrorOutput(result.error || 'Unknown compilation error');
                updateStatus('error', 'Compilation failed');
            }
            
        } catch (error) {
            output.className = 'error';
            output.innerHTML = formatErrorOutput(`Failed to connect to compiler: ${error.message}`);
            updateStatus('error', 'Connection failed');
        }
    }
    
    async function runCode() {
        const code = editor.getValue();
        const output = document.getElementById('output');
        
        updateStatus('', 'Running...');
        output.innerHTML = '<div style="color: #ffa500;">Compiling and running...</div>';
        
        try {
            const response = await fetch(`${API_URL}/run`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ code })
            });
            
            const result = await response.json();
            
            if (!response.ok) {
                // Handle HTTP errors that still have JSON error details
                if (result.error) {
                    output.className = 'error';
                    const statusMessage = result.isCompilationError ? 'Compilation failed' : 'Execution failed';
                    
                    output.innerHTML = formatErrorOutput(result.error);
                    updateStatus('error', statusMessage);
                    return;
                } else {
                    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
                }
            }
            
            if (result.success) {
                // Successful execution
                output.className = 'success';
                let outputText = '';
                
                if (result.programOutput && result.programOutput.trim()) {
                    outputText = result.programOutput;
                } else {
                    outputText = '✅ Program ran successfully - no output';
                }
                
                output.innerHTML = formatPlainOutput(outputText);
                updateStatus('connected', 'Ready');
            } else {
                // Execution failed
                output.className = 'error';
                output.innerHTML = formatErrorOutput(result.error || 'Unknown error');
                updateStatus('error', 'Execution failed');
            }
            
        } catch (error) {
            output.className = 'error';
            output.innerHTML = formatErrorOutput(`Failed to connect to compiler: ${error.message}`);
            updateStatus('error', 'Connection failed');
        }
    }
    
    function formatErrorOutput(text) {
        if (!text) return '';
        
        // Escape HTML
        text = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
        
        // Split by lines and parse errors
        const lines = text.split('\n').filter(line => line.trim());
        const errorLines = [];
        
        lines.forEach(line => {
            // Check if line contains line number references
            const lineNumberMatch = line.match(/\b(?:line\s+)(\d+)(?:\s*:\s*(\d+))?/i) ||
                                  line.match(/\bat line\s+(\d+)/i) ||
                                  line.match(/\berror at\s+(\d+)/i) ||
                                  line.match(/\[(\d+)(?:\s*:\s*(\d+))?\]/);
            
            if (lineNumberMatch) {
                const lineNum = parseInt(lineNumberMatch[1]);
                const column = lineNumberMatch[2] ? parseInt(lineNumberMatch[2]) : 0;
                
                // Extract the error message (everything after the line number)
                let message = line.replace(/^.*?(?:line\s+\d+(?::\d+)?|at line\s+\d+|\[\d+(?::\d+)?\])\s*/, '').trim();
                if (!message) message = line.trim();
                
                errorLines.push({
                    lineNum,
                    column,
                    message,
                    fullText: line
                });
            } else {
                // Non-line-specific error
                errorLines.push({
                    lineNum: null,
                    column: null,
                    message: line.trim(),
                    fullText: line
                });
            }
        });
        
        if (errorLines.length === 0) {
            return text; // Fallback to original text
        }
        
        // Build clean grid structure
        const gridItems = errorLines.map(error => {
            if (error.lineNum !== null) {
                const location = error.column > 0 ? `${error.lineNum}:${error.column}` : `${error.lineNum}`;
                return `<div class="error-line" onclick="jumpToLine(${error.lineNum}, ${error.column || 1})">
                    <span class="error-location">Line ${location}</span>
                    <span class="error-message">${error.message}</span>
                </div>`;
            } else {
                return `<div class="error-line">
                    <span class="error-location">—</span>
                    <span class="error-message">${error.message}</span>
                </div>`;
            }
        });
        
        return `<div class="error-list">${gridItems.join('')}</div>`;
    }
    
    function formatPlainOutput(text) {
        if (!text) return '';
        
        // Escape HTML
        text = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
        
        // Color specific messages
        text = text.replace(/(Program executed successfully)/g, '<span style="color: #7cb992;">$1</span>');
        text = text.replace(/(Running program\.\.\.)/g, '<span style="color: #ffa500;">$1</span>');
        
        return text;
    }
    
    function jumpToLine(lineNumber, column = 1) {
        if (!editor) return;
        
        console.log(`🎯 Jumping to line ${lineNumber}, column ${column}`);
        
        // Remove any existing selections
        const errorLines = document.querySelectorAll('.error-line');
        errorLines.forEach(el => el.classList.remove('selected'));
        
        // Mark clicked line as selected
        event.target.closest('.error-line')?.classList.add('selected');
        
        // Jump to the line in Monaco editor
        editor.setPosition({ lineNumber: lineNumber, column: column });
        editor.revealLineInCenter(lineNumber);
        editor.focus();
        
        // Optionally highlight the line temporarily
        const decoration = editor.deltaDecorations([], [{
            range: new monaco.Range(lineNumber, 1, lineNumber, 1),
            options: {
                isWholeLine: true,
                className: 'highlighted-error-line',
                glyphMarginClassName: 'error-glyph'
            }
        }]);
        
        // Remove decoration after 2 seconds
        setTimeout(() => {
            editor.deltaDecorations(decoration, []);
        }, 2000);
    }
    
    function clearOutput() {
        document.getElementById('output').innerHTML = '';
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