// Page-health checks for every template at every breakpoint:
// no JS errors, no failed requests, no horizontal overflow, dark theme intact,
// site chrome present, and the logo image resolving.
const { test, expect } = require("@playwright/test");
const { PAGES, VIEWPORTS, collectProblems } = require("./helpers");

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
