const fs = require('fs');
const path = require('path');

const specSourcePath = path.resolve(__dirname, '../../compiler/spec.md');
const specDestPath = path.resolve(__dirname, '../src/spec.md');

try {
  // Read the spec file and add front matter for proper rendering
  const specContent = fs.readFileSync(specSourcePath, 'utf8');
  
  // Add front matter to make it render as a blog-style page
  const frontMatter = `---
layout: page
title: "Osprey Language Specification"
description: "Complete language specification and syntax reference for the Osprey programming language"
date: ${new Date().toISOString().split('T')[0]}
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/"
---

`;
  
  const contentWithFrontMatter = frontMatter + specContent;
  fs.writeFileSync(specDestPath, contentWithFrontMatter, 'utf8');
  
  console.log('✅ Copied spec.md from compiler directory');
} catch (error) {
  console.error('❌ Failed to copy spec.md:', error.message);
  process.exit(1);
} 