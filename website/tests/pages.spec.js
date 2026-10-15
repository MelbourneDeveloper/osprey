// Page-health checks for every template at every breakpoint:
// no JS errors, no failed requests, no horizontal overflow, dark theme intact,
// site chrome present, the logo image resolving, and Osprey code highlighted.
const fs = require("node:fs");
const path = require("node:path");
const { test, expect } = require("@playwright/test");
const { PAGES, VIEWPORTS, collectProblems, paintedColours } = require("./helpers");

const SITE = path.resolve(__dirname, "../_site");
const OSPREY_CODE = /<code([^>]*)class="language-osprey(?:-ml)?"[^>]*>([\s\S]*?)<\/code>/g;
const COMMENT_TOKEN = /<span class="token comment">[\s\S]*?<\/span>/g;
const OSPREY_BLOCKS = "pre > code.language-osprey, pre > code.language-osprey-ml";

// A block is unhighlighted when it has no tokens or an ML `(* … *)` comment
// outside a comment token. The /wasm/ studio fills #src-code at runtime; the
// browser check below waits for its tokens.
const unhighlightedBlocks = (file) =>
  [...fs.readFileSync(path.join(SITE, file), "utf8").matchAll(OSPREY_CODE)]
    .filter(([, attributes, body]) => !attributes.includes('id="src-code"'))
    .filter(([, , body]) => !body.includes('class="token') || body.replace(COMMENT_TOKEN, "").replace(/<[^>]+>/g, "").includes("(*"))
    .map(([, , body]) => `${file}: ${body.slice(0, 60)}`);

test("every built page ships its Osprey code highlighted", () => {
  const pages = fs.readdirSync(SITE, { recursive: true }).filter((file) => file.endsWith(".html"));
  expect(pages).toContain("index.html");
  expect(pages.flatMap(unhighlightedBlocks)).toEqual([]);
});

const DARK_BACKGROUNDS = ["rgb(12, 19, 37)", "rgb(7, 13, 31)", "rgb(25, 31, 50)"];

for (const pg of PAGES) {
  test.describe(`page: ${pg.name}`, () => {
    test(`loads cleanly with no errors (${pg.path})`, async ({ page }) => {
      const problems = collectProblems(page, pg.name);
      const resp = await page.goto(pg.path, { waitUntil: "networkidle" });
      expect(resp.status(), "HTTP status").toBeLessThan(400);
      await page.waitForTimeout(300);

      expect(problems.pageErrors, "uncaught page errors").toEqual([]);
      expect(problems.consoleErrors, "console errors").toEqual([]);
      expect(problems.failedRequests, "failed requests").toEqual([]);
    });

    test(`has dark theme + header + footer (${pg.path})`, async ({ page }) => {
      await page.goto(pg.path, { waitUntil: "domcontentloaded" });
      const bg = await page.evaluate(() => getComputedStyle(document.body).backgroundColor);
      expect(DARK_BACKGROUNDS, `body background was ${bg}`).toContain(bg);

      await expect(page.locator(".site-header")).toBeVisible();
      // The theme toggle must stay hidden — the design is dark-only.
      await expect(page.locator(".theme-toggle")).toBeHidden();
      // Full-screen app-like pages (playground, wasm studio) hide the footer.
      if (pg.kind !== "app" && !pg.fullscreen) {
        await expect(page.locator(".site-footer")).toBeVisible();
      }
    });

    test(`logo image resolves (${pg.path})`, async ({ page }) => {
      await page.goto(pg.path, { waitUntil: "domcontentloaded" });
      const bgImage = await page.evaluate(() =>
        getComputedStyle(document.querySelector(".logo"), "::before").backgroundImage
      );
      const match = bgImage.match(/url\(["']?([^"')]+)["']?\)/);
      expect(match, "logo ::before background-image url").toBeTruthy();
      const res = await page.request.get(match[1]);
      expect(res.status(), "logo request status").toBe(200);
    });

    test(`Osprey code paints its tokens (${pg.path})`, async ({ page }) => {
      await page.goto(pg.path, { waitUntil: "domcontentloaded" });
      for (const block of await page.locator(OSPREY_BLOCKS).all()) {
        await expect(block.locator(".token").first()).toBeAttached();
        expect(await paintedColours(block.or(block.locator(".token"))), await block.textContent()).toBeGreaterThan(1);
      }
    });

    // The playground is a fixed-height Monaco app — horizontal overflow there
    // is owned by the editor, not the site layout.
    if (pg.kind !== "app") {
      for (const vp of VIEWPORTS) {
        test(`no horizontal overflow @ ${vp.name} (${pg.path})`, async ({ page }) => {
          await page.setViewportSize({ width: vp.width, height: vp.height });
          await page.goto(pg.path, { waitUntil: "networkidle" });
          await page.waitForTimeout(200);
          const overflow = await page.evaluate(
            () => document.documentElement.scrollWidth - window.innerWidth
          );
          expect(overflow, `document overflows by ${overflow}px`).toBeLessThanOrEqual(2);
        });
      }
    }
  });
}
