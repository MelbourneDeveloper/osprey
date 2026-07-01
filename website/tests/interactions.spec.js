// User-interaction tests with explicit assertions: navigation, the hero,
// the mobile menu drawer, and the prose Table-of-Contents.
const { test, expect } = require("@playwright/test");

test.describe("desktop interactions", () => {
  test.use({ viewport: { width: 1440, height: 900 } });

  test("nav links navigate to the right pages", async ({ page }) => {
    await page.goto("/");
    await page.locator(".nav-link", { hasText: "Spec" }).click();
    await expect(page).toHaveURL(/\/spec\/$/);
    await page.locator(".nav-link", { hasText: "Docs" }).click();
    await expect(page).toHaveURL(/\/docs\/$/);
  });

  test("logo returns to home", async ({ page }) => {
    await page.goto("/spec/0001-introduction/");
    await page.locator(".logo").click();
    await expect(page).toHaveURL(/\/$/);
    await expect(page.locator(".hero")).toBeVisible();
  });

  test("hero fills the viewport before scroll", async ({ page }) => {
    await page.goto("/");
    const { heroH, winH } = await page.evaluate(() => ({
      heroH: document.querySelector(".hero").getBoundingClientRect().height,
      winH: window.innerHeight,
    }));
    // Hero + the in-flow sticky header should cover the viewport.
    expect(heroH).toBeGreaterThanOrEqual(winH - 70);
    // The next section must start below the fold.
    const nextTop = await page.evaluate(
      () => document.querySelector(".blue-box-installation").getBoundingClientRect().top
    );
    expect(nextTop).toBeGreaterThanOrEqual(winH - 5);
  });

  test("hero has a title and two working CTAs", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator(".hero-title")).toContainText("Osprey");
    await expect(page.locator(".hero-actions .btn")).toHaveCount(2);
    await page.locator(".hero-actions .btn", { hasText: "Try Osprey Online" }).click();
    await expect(page).toHaveURL(/\/playground\/$/);
  });

  test("wasm studio serves generated modules and toggles flavors", async ({ page }) => {
    for (const asset of [
      "/wasm/wasi-shim.mjs",
      "/wasm/studio.osp",
      "/wasm/studio.ospml",
      "/wasm/build/studio.osp.wasm",
      "/wasm/build/studio.ospml.wasm",
    ]) {
      const res = await page.request.get(asset);
      expect(res.status(), `${asset} status`).toBe(200);
    }

    await page.goto("/wasm/");
    await expect(page.locator("#metrics .metric")).toHaveCount(5, { timeout: 15_000 });
    await expect(page.locator("#flavor-bytes")).toContainText("KB wasm");
    await page.locator("#flavor-ospml").click();
    await expect(page.locator("#flavor-name")).toContainText("studio.ospml");
    await expect(page.locator("#metrics .metric")).toHaveCount(5);
  });

  test("real-world example code is not clipped", async ({ page }) => {
    await page.goto("/");
    const clips = await page.evaluate(() =>
      [...document.querySelectorAll(".showcase-grid .card-code pre")].map(
        (pre) => pre.scrollWidth - pre.clientWidth
      )
    );
    expect(clips.length).toBeGreaterThan(0);
    for (const c of clips) expect(c, "code block horizontal clip (px)").toBeLessThanOrEqual(2);
  });

  test("prose page shows a TOC with scroll-spy", async ({ page }) => {
    await page.goto("/spec/0001-introduction/");
    const toc = page.locator(".toc");
    await expect(toc).toBeVisible();
    const links = page.locator(".toc-link");
    expect(await links.count()).toBeGreaterThan(0);

    // Clicking a TOC link jumps to that section.
    const second = links.nth(1);
    const href = await second.getAttribute("href");
    await second.click();
    await expect(page).toHaveURL(new RegExp(href.replace(/[.*+?^${}()|[\]\\]/g, "\\$&") + "$"));
    const targetVisible = await page.locator(href).isVisible();
    expect(targetVisible).toBeTruthy();

    // Scroll spy marks an active link.
    await page.evaluate(() => window.scrollBy(0, 1200));
    await page.waitForTimeout(400);
    await expect(page.locator(".toc-link.active")).toHaveCount(1);
  });

  test("prose headings render as plain text, not links", async ({ page }) => {
    await page.goto("/spec/0001-introduction/");
    const h2 = page.locator(".prose h2").first();
    const color = await h2.evaluate((el) => getComputedStyle(el).color);
    // on-surface (#dce1fb) — NOT the cyan link colour (#77d7f4 / #bdeeff).
    expect(color).toBe("rgb(220, 225, 251)");
  });
});

test.describe("mobile interactions", () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test("mobile menu opens and closes", async ({ page }) => {
    await page.goto("/");
    const toggle = page.locator("#mobile-menu-toggle");
    const links = page.locator(".nav-links");
    await expect(toggle).toBeVisible();
    await expect(links).toBeHidden();
    await toggle.click();
    await expect(links).toBeVisible();
    await expect(page.locator(".nav-link", { hasText: "Docs" })).toBeVisible();
    await toggle.click();
    await expect(links).toBeHidden();
  });

  test("mobile menu link navigates", async ({ page }) => {
    await page.goto("/");
    await page.locator("#mobile-menu-toggle").click();
    await page.locator(".nav-link", { hasText: "Blog" }).click();
    await expect(page).toHaveURL(/\/blog\/$/);
  });

  test("hero actions stack full-width on mobile", async ({ page }) => {
    await page.goto("/");
    const widths = await page.evaluate(() => {
      const btns = [...document.querySelectorAll(".hero-actions .btn")];
      return { btn: btns[0].getBoundingClientRect().width, actions: document.querySelector(".hero-actions").getBoundingClientRect().width };
    });
    // Each button should span (nearly) the full actions row.
    expect(widths.btn).toBeGreaterThan(widths.actions * 0.9);
  });

  test("TOC sidebar is hidden on mobile", async ({ page }) => {
    await page.goto("/spec/0001-introduction/");
    await expect(page.locator(".toc-aside")).toBeHidden();
  });
});
