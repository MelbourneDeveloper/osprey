const syntaxHighlight = require("@11ty/eleventy-plugin-syntaxhighlight");
const markdownIt = require("markdown-it");
const markdownItAnchor = require("markdown-it-anchor");
const { DateTime } = require("luxon");

module.exports = function (eleventyConfig) {
  // Add plugins with enhanced configuration
  eleventyConfig.addPlugin(syntaxHighlight, {
    templateFormats: ["*"],
    preAttributes: {
      tabindex: 0,
      // Add data attributes for enhanced styling
      "data-language": function ({ language, content, options }) {
        return language;
      }
    },
    codeAttributes: {},
    init: function ({ Prism }) {
      // Define Osprey language syntax for Prism
      Prism.languages.osprey = {
        'comment': [
          {
            pattern: /(^|[^\\])\/\*[\s\S]*?(?:\*\/|$)/,
            lookbehind: true
          },
          {
            pattern: /(^|[^\\:])\/\/.*/,
            lookbehind: true
          }
        ],
        'string': {
          pattern: /"(?:[^"\\]|\\.)*"/,
          greedy: true
        },
        'interpolation': {
          pattern: /\$\{[^}]+\}/,
          inside: {
            'punctuation': /^\$\{|\}$/,
            'expression': {
              pattern: /[\s\S]+/,
              inside: Prism.languages.osprey || {}
            }
          }
        },
        'keyword': /\b(?:fn|let|match|type|if|else|while|for|return|module|import|export|true|false|and|or|not|where|await|async|fiber|Some|None|Ok|Err|Result|Option)\b/,
        'type': /\b(?:Int|String|Bool|Float|List|Map|Set|Result|Option|Error|Json)\b/,
        'function': /\b[a-zA-Z_][a-zA-Z0-9_]*(?=\s*\()/,
        'number': /\b(?:0x[\da-f]+|\d*\.?\d+(?:e[+-]?\d+)?)\b/i,
        'operator': /\+\+|--|&&|\|\||==|!=|<=|>=|<|>|\+|-|\*|\/|%|=|\|>|=>|->|<-/,
        'punctuation': /[{}[\];(),.:]/,
        'pipe': /\|>/,
        'pattern': /\b[A-Z][a-zA-Z0-9_]*(?=\s*\{)/,
        'variable': /\b[a-z_][a-zA-Z0-9_]*\b/
      };
    }
  });

  // Transform to automatically highlight inline code blocks
  eleventyConfig.addTransform("highlight-code", function (content, outputPath) {
    // Only process HTML files
    if (outputPath && outputPath.endsWith(".html")) {
      // Replace raw code blocks with properly highlighted ones
      content = content.replace(
        /<pre class="language-osprey"><code class="language-osprey">([\s\S]*?)<\/code><\/pre>/g,
        function (match, code) {
          // Decode HTML entities and clean up the code
          let decodedCode = code
            .replace(/&lt;/g, '<')
            .replace(/&gt;/g, '>')
            .replace(/&amp;/g, '&')
            .replace(/&quot;/g, '"')
            .replace(/&#39;/g, "'")
            // Remove any HTML tags that got inserted into the code
            .replace(/<\/?[^>]+(>|$)/g, "")
            .trim();

          // Use Prism to highlight the code
          const Prism = require('prismjs');

          // Load Osprey language if not already loaded
          if (!Prism.languages.osprey) {
            Prism.languages.osprey = {
              'comment': [
                {
                  pattern: /(^|[^\\])\/\*[\s\S]*?(?:\*\/|$)/,
                  lookbehind: true
                },
                {
                  pattern: /(^|[^\\:])\/\/.*/,
                  lookbehind: true
                }
              ],
              'string': {
                pattern: /"(?:[^"\\]|\\.)*"/,
                greedy: true
              },
              'interpolation': {
                pattern: /\$\{[^}]+\}/,
                inside: {
                  'punctuation': /^\$\{|\}$/,
                  'expression': {
                    pattern: /[\s\S]+/,
                    inside: Prism.languages.osprey || {}
                  }
                }
              },
              'keyword': /\b(?:fn|let|match|type|if|else|while|for|return|module|import|export|true|false|and|or|not|where|await|async|fiber|Some|None|Ok|Err|Result|Option)\b/,
              'type': /\b(?:Int|String|Bool|Float|List|Map|Set|Result|Option|Error|Json)\b/,
              'function': /\b[a-zA-Z_][a-zA-Z0-9_]*(?=\s*\()/,
              'number': /\b(?:0x[\da-f]+|\d*\.?\d+(?:e[+-]?\d+)?)\b/i,
              'operator': /\+\+|--|&&|\|\||==|!=|<=|>=|<|>|\+|-|\*|\/|%|=|\|>|=>|->|<-/,
              'punctuation': /[{}[\];(),.:]/,
              'pipe': /\|>/,
              'pattern': /\b[A-Z][a-zA-Z0-9_]*(?=\s*\{)/,
              'variable': /\b[a-z_][a-zA-Z0-9_]*\b/
            };
          }

          const highlightedCode = Prism.highlight(decodedCode, Prism.languages.osprey, 'osprey');
          return `<pre class="language-osprey" tabindex="0" data-language="osprey"><code class="language-osprey">${highlightedCode}</code></pre>`;
        }
      );
    }
    return content;
  });

  // Pass through files
  eleventyConfig.addPassthroughCopy("src/assets");
  eleventyConfig.addPassthroughCopy("src/css");
  eleventyConfig.addPassthroughCopy("src/js");
  eleventyConfig.addPassthroughCopy("src/playground");

  // Watch targets
  eleventyConfig.addWatchTarget("src/css/");
  eleventyConfig.addWatchTarget("src/js/");
  eleventyConfig.addWatchTarget("../compiler/spec.md");

  // Filters
  eleventyConfig.addFilter("readableDate", dateObj => {
    return DateTime.fromJSDate(dateObj, { zone: 'utc' }).toFormat("dd LLL yyyy");
  });

  eleventyConfig.addFilter('htmlDateString', (dateObj) => {
    return DateTime.fromJSDate(dateObj, { zone: 'utc' }).toFormat('yyyy-LL-dd');
  });

  // Layout aliases - make all markdown files use blog layout by default
  eleventyConfig.addLayoutAlias('default', 'base.njk');
  eleventyConfig.addLayoutAlias('docs', 'base.njk');

  // Set default layout for all markdown files
  eleventyConfig.addGlobalData("layout", "base");

  // Collections
  eleventyConfig.addCollection("blog", function (collectionApi) {
    return collectionApi.getFilteredByGlob("src/blog/**/*.md")
      .filter(post => !post.inputPath.includes('/index.md'))
      .sort((a, b) => {
        return b.date - a.date;
      });
  });

  eleventyConfig.addCollection("docs", function (collectionApi) {
    return collectionApi.getFilteredByGlob("src/docs/**/*.md");
  });

  eleventyConfig.addCollection("howto", function (collectionApi) {
    return collectionApi.getFilteredByGlob("src/docs/how-to/**/*.md");
  });

  eleventyConfig.addCollection("reference", function (collectionApi) {
    return collectionApi.getFilteredByGlob("src/docs/**/*.md");
  });

  // Markdown configuration - FIXED to properly handle HTML
  let markdownOptions = {
    html: true,
    breaks: false,  // Don't convert line breaks to <br>
    linkify: true,
    typographer: true
  };

  let markdownLib = markdownIt(markdownOptions).use(markdownItAnchor, {
    permalink: markdownItAnchor.permalink.ariaHidden({
      placement: "after",
      class: "header-anchor",
      symbol: "#",
      ariaHidden: false,
    }),
    level: [2, 3, 4],
    slugify: eleventyConfig.getFilter("slug")
  });

  // Remove the custom paragraph rules that were causing issues
  // Let markdown-it handle paragraphs normally

  eleventyConfig.setLibrary("md", markdownLib);

  // Short codes
  eleventyConfig.addShortcode("year", () => `${new Date().getFullYear()}`);

  // Enhanced code example shortcode for Osprey with proper syntax highlighting
  eleventyConfig.addPairedShortcode("osprey", function (content, title = "") {
    const titleHtml = title ? `<div class="code-title">${title}</div>` : "";

    // Use the same highlighting as the plugin
    const Prism = require('prismjs');

    // Ensure Osprey language is loaded
    if (!Prism.languages.osprey) {
      Prism.languages.osprey = {
        'comment': [
          {
            pattern: /(^|[^\\])\/\*[\s\S]*?(?:\*\/|$)/,
            lookbehind: true
          },
          {
            pattern: /(^|[^\\:])\/\/.*/,
            lookbehind: true
          }
        ],
        'string': {
          pattern: /"(?:[^"\\]|\\.)*"/,
          greedy: true
        },
        'interpolation': {
          pattern: /\$\{[^}]+\}/,
          inside: {
            'punctuation': /^\$\{|\}$/,
            'expression': {
              pattern: /[\s\S]+/,
              inside: Prism.languages.osprey || {}
            }
          }
        },
        'keyword': /\b(?:fn|let|match|type|if|else|while|for|return|module|import|export|true|false|and|or|not|where|await|async|fiber|Some|None|Ok|Err|Result|Option)\b/,
        'type': /\b(?:Int|String|Bool|Float|List|Map|Set|Result|Option|Error|Json)\b/,
        'function': /\b[a-zA-Z_][a-zA-Z0-9_]*(?=\s*\()/,
        'number': /\b(?:0x[\da-f]+|\d*\.?\d+(?:e[+-]?\d+)?)\b/i,
        'operator': /\+\+|--|&&|\|\||==|!=|<=|>=|<|>|\+|-|\*|\/|%|=|\|>|=>|->|<-/,
        'punctuation': /[{}[\];(),.:]/,
        'pipe': /\|>/,
        'pattern': /\b[A-Z][a-zA-Z0-9_]*(?=\s*\{)/,
        'variable': /\b[a-z_][a-zA-Z0-9_]*\b/
      };
    }

    const highlightedCode = Prism.highlight(content.trim(), Prism.languages.osprey, 'osprey');

    return `<div class="code-example">
      ${titleHtml}
      <pre class="language-osprey" tabindex="0" data-language="osprey"><code class="language-osprey">${highlightedCode}</code></pre>
    </div>`;
  });

  // Playground embed shortcode
  eleventyConfig.addShortcode("playground", function (height = "600px", code = "") {
    const encodedCode = code ? encodeURIComponent(code) : "";
    return `<div class="playground-embed" style="width: 100%; height: ${height}; border: 1px solid var(--color-neutral-300); border-radius: var(--radius-lg); overflow: hidden;">
      <iframe 
        src="/playground/${encodedCode ? `?code=${encodedCode}` : ''}" 
        style="width: 100%; height: 100%; border: none;"
        loading="lazy"
        allow="clipboard-write"
        title="Osprey Language Playground">
      </iframe>
    </div>`;
  });

  // Interactive code example shortcode
  eleventyConfig.addPairedShortcode("interactive", function (content, title = "") {
    const encodedCode = encodeURIComponent(content.trim());
    return `<div class="interactive-example">
      ${title ? `<div class="example-title">${title}</div>` : ""}
      <div class="playground-embed" style="width: 100%; height: 400px; border-radius: 8px; overflow: hidden;">
        <iframe 
          src="/playground/embed.html#${encodedCode}" 
          style="width: 100%; height: 100%; border: none;"
          loading="lazy"
          allow="clipboard-write"
          title="${title || 'Interactive Osprey Example'}">
        </iframe>
      </div>
    </div>`;
  });

  return {
    dir: {
      input: "src",
      output: "_site",
      includes: "_includes",
      layouts: "_layouts",
      data: "_data"
    },
    markdownTemplateEngine: "njk",
    htmlTemplateEngine: "njk",
    dataTemplateEngine: "njk",
    // Server configuration for dev container compatibility
    serverOptions: {
      port: 8080,
      host: "0.0.0.0",
      showAllHosts: true
    }
  };
}; 