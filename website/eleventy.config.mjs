// Eleventy config — Osprey website built on the eleventy-plugin-techdoc theme.
// The theme (a structural-only "virtual theme") provides the HTML shell, head
// SEO/JSON-LD, nav/footer, dark mode, and auto-generates sitemap/robots/feed/
// llms.txt. It also registers syntaxhighlight, rss, navigation, markdown and a
// `year` shortcode — so this config must NOT re-register those. We add only
// what is Osprey-specific: the Prism grammar for `.osp`, a transform that
// highlights raw `language-osprey` blocks, and the playground shortcodes.
import techdoc from "eleventy-plugin-techdoc";
import Prism from "prismjs";
import { DateTime } from "luxon";

// Osprey Prism grammar — shared by the syntaxhighlight plugin and the transform.
const ospreyGrammar = {
  comment: [
    { pattern: /(^|[^\\])\/\*[\s\S]*?(?:\*\/|$)/, lookbehind: true },
    { pattern: /(^|[^\\:])\/\/.*/, lookbehind: true },
  ],
  string: { pattern: /"(?:[^"\\]|\\.)*"/, greedy: true },
  interpolation: {
    pattern: /\$\{[^}]+\}/,
    inside: { punctuation: /^\$\{|\}$/ },
  },
  keyword:
    /\b(?:fn|let|mut|match|type|effect|perform|handle|in|extern|spawn|await|yield|if|else|import|module|true|false|where|Unit|Result|Option|Some|None|Ok|Err)\b/,
  type: /\b(?:int|float|string|bool|List|Map|Set|Ptr|Channel|Fiber|Json|HttpResponse)\b/,
  function: /\b[a-zA-Z_][a-zA-Z0-9_]*(?=\s*\()/,
  number: /\b(?:0x[\da-f]+|\d*\.?\d+(?:e[+-]?\d+)?)\b/i,
  operator: /\|>|->|=>|<-|\+|-|\*|\/|%|==|!=|<=|>=|<|>|=|!|&&|\|\|/,
  punctuation: /[{}[\];(),.:]/,
};

function ensureOsprey() {
  if (!Prism.languages.osprey) Prism.languages.osprey = ospreyGrammar;
}

export default function (eleventyConfig) {
  eleventyConfig.addPlugin(techdoc, {
    site: {
      name: "Osprey",
      url: "https://ospreylang.dev",
      description:
        "A modern functional language with compile-time effect safety, lightweight fiber concurrency, and immutable persistent collections.",
    },
    // Keep the existing blog index + docs pages; only adopt the theme's shell,
    // SEO and generated sitemap/robots/feed/llms.txt. (New designs land later.)
    features: { blog: false, docs: false, darkMode: true, i18n: false },
    i18n: { defaultLanguage: "en", languages: ["en"] },
  });

  // Register the Osprey grammar so the theme's bundled syntaxhighlight (and the
  // transform below) can colour `.osp` snippets.
  ensureOsprey();

  // Highlight raw `<pre class="language-osprey">` blocks that ship as literal
  // HTML in the marketing pages (not processed by the markdown highlighter).
  eleventyConfig.addTransform("osprey-highlight", function (content, outputPath) {
    if (!outputPath || !outputPath.endsWith(".html")) return content;
    ensureOsprey();
    return content.replace(
      /<pre class="language-osprey"><code class="language-osprey">([\s\S]*?)<\/code><\/pre>/g,
      (_m, code) => {
        const decoded = code
          .replace(/&lt;/g, "<")
          .replace(/&gt;/g, ">")
          .replace(/&amp;/g, "&")
          .replace(/&quot;/g, '"')
          .replace(/&#39;/g, "'")
          .replace(/<\/?[^>]+(>|$)/g, "")
          .trim();
        const html = Prism.highlight(decoded, Prism.languages.osprey, "osprey");
        return `<pre class="language-osprey" tabindex="0" data-language="osprey"><code class="language-osprey">${html}</code></pre>`;
      }
    );
  });

  // Playground embed shortcode (used by docs/blog markdown).
  eleventyConfig.addPairedShortcode("interactive", function (content, title = "") {
    const encoded = encodeURIComponent(content.trim());
    return `<div class="interactive-example">${
      title ? `<div class="example-title">${title}</div>` : ""
    }<div class="playground-embed"><iframe src="/playground/#${encoded}" loading="lazy" allow="clipboard-write" title="${
      title || "Interactive Osprey Example"
    }"></iframe></div></div>`;
  });

  // Osprey's own CSS, JS and the Monaco-based playground ship as static assets.
  eleventyConfig.addPassthroughCopy("src/assets");
  eleventyConfig.addPassthroughCopy("src/css");
  eleventyConfig.addPassthroughCopy("src/js");
  eleventyConfig.addPassthroughCopy("src/playground");
  // Publish WebAssembly demo assets for the native /wasm/ page. The deploy
  // pipeline runs `make wasm-site` first so generated binaries land here.
  eleventyConfig.addPassthroughCopy({
    "../examples/wasm/build/studio.osp.wasm": "wasm/build/studio.osp.wasm",
  });
  eleventyConfig.addPassthroughCopy({
    "../examples/wasm/build/studio.ospml.wasm": "wasm/build/studio.ospml.wasm",
  });
  eleventyConfig.addPassthroughCopy({ "../examples/wasm/wasi-shim.mjs": "wasm/wasi-shim.mjs" });
  eleventyConfig.addPassthroughCopy({ "../examples/wasm/studio.osp": "wasm/studio.osp" });
  eleventyConfig.addPassthroughCopy({ "../examples/wasm/studio.ospml": "wasm/studio.ospml" });

  eleventyConfig.addWatchTarget("src/css/");
  eleventyConfig.addWatchTarget("src/js/");
  eleventyConfig.addWatchTarget("../examples/wasm/");

  // Map the site's existing layout names onto the theme's base layout. Existing
  // pages declare `layout: page`, `layout: page.njk` or `layout: base.njk`; the
  // theme ships `layouts/base.njk`. Aliasing avoids touching every page.
  eleventyConfig.addLayoutAlias("base", "layouts/base.njk");
  eleventyConfig.addLayoutAlias("base.njk", "layouts/base.njk");
  // Long-form pages (docs, spec, blog posts, status) share ONE prose design.
  eleventyConfig.addLayoutAlias("page", "layouts/prose.njk");
  eleventyConfig.addLayoutAlias("page.njk", "layouts/prose.njk");

  // The blog index renders this collection (theme blog auto-pages are disabled).
  eleventyConfig.addCollection("blog", (api) =>
    api
      .getFilteredByGlob("src/blog/**/*.md")
      .filter((p) => !p.inputPath.includes("/index."))
      .sort((a, b) => b.date - a.date)
  );

  // Date filters the blog listing uses (the theme exposes dateFormat/isoDate).
  eleventyConfig.addFilter("readableDate", (d) =>
    DateTime.fromJSDate(d, { zone: "utc" }).toFormat("dd LLL yyyy")
  );
  eleventyConfig.addFilter("htmlDateString", (d) =>
    DateTime.fromJSDate(d, { zone: "utc" }).toFormat("yyyy-LL-dd")
  );

  return {
    dir: { input: "src", output: "_site", data: "_data" },
    markdownTemplateEngine: "njk",
    htmlTemplateEngine: "njk",
  };
}
