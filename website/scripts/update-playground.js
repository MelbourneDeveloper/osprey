#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Paths
const MEGA_SHOWCASE_PATH = path.join(__dirname, '../../compiler/examples/tested/basics/osprey_mega_showcase.osp');
const PLAYGROUND_PATH = path.join(__dirname, '../src/playground/index.md');

function escapeBackticks(str) {
    return str.replace(/`/g, '\\`');
}

function updatePlayground() {
    try {
        // Read the mega showcase file
        if (!fs.existsSync(MEGA_SHOWCASE_PATH)) {
            console.error('❌ Mega showcase file not found:', MEGA_SHOWCASE_PATH);
            return false;
        }
        
        const megaShowcaseCode = fs.readFileSync(MEGA_SHOWCASE_PATH, 'utf8');
        console.log('📖 Read mega showcase code (' + megaShowcaseCode.length + ' chars)');
        
        // Read the current playground file
        if (!fs.existsSync(PLAYGROUND_PATH)) {
            console.error('❌ Playground file not found:', PLAYGROUND_PATH);
            return false;
        }
        
        const playgroundContent = fs.readFileSync(PLAYGROUND_PATH, 'utf8');
        
        // Find the editor value section and replace it
        const valueStartMarker = "value: `";
        const valueEndMarker = "`,\n            language: 'osprey'";
        
        const startIndex = playgroundContent.indexOf(valueStartMarker);
        const endIndex = playgroundContent.indexOf(valueEndMarker);
        
        if (startIndex === -1 || endIndex === -1) {
            console.error('❌ Could not find editor value section in playground file');
            return false;
        }
        
        // Escape backticks in the mega showcase code
        const escapedCode = escapeBackticks(megaShowcaseCode);
        
        // Replace the editor value
        const newPlaygroundContent = 
            playgroundContent.substring(0, startIndex + valueStartMarker.length) +
            escapedCode +
            playgroundContent.substring(endIndex);
        
        // Write the updated playground file
        fs.writeFileSync(PLAYGROUND_PATH, newPlaygroundContent, 'utf8');
        
        console.log('✅ Successfully updated playground with mega showcase example!');
        console.log('🎯 Playground now has the latest comprehensive sandboxable features demo');
        
        return true;
        
    } catch (error) {
        console.error('❌ Error updating playground:', error.message);
        return false;
    }
}

// Run the update
if (require.main === module) {
    console.log('🚀 Updating Osprey Playground with Mega Showcase...');
    const success = updatePlayground();
    process.exit(success ? 0 : 1);
}

module.exports = { updatePlayground }; 