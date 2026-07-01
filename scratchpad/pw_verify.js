const { chromium } = require('playwright');

const BASE = 'http://localhost:8099';
const PAGES = [
  '/docs/functions/channel/',
  '/docs/types/any/',
  '/docs/keywords/match/',
  '/docs/functions/map/',
  '/docs/keywords/fn/',
];

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage({ viewport: { width: 1400, height: 1600 } });
  let allOk = true;
  for (const path of PAGES) {
    await page.goto(BASE + path, { waitUntil: 'networkidle' });
    // count flavor badges via the data-flavor attribute the transform sets
    const flavors = await page.$$eval('pre[data-flavor]', els =>
      els.map(e => e.getAttribute('data-flavor')));
    const hasDefault = flavors.includes('default');
    const hasMl = flavors.includes('ml');
    const ok = hasDefault && hasMl;
    allOk = allOk && ok;
    const name = path.replace(/\//g, '_').replace(/^_|_$/g, '');
    await page.screenshot({ path: `scratchpad/shot_${name}.png`, fullPage: true });
    console.log(`${ok ? 'PASS' : 'FAIL'}  ${path}  badges=[${flavors.join(', ')}]`);
  }
  await browser.close();
  console.log(allOk ? '\nALL PAGES SHOW BOTH FLAVORS ✓' : '\nSOME PAGES MISSING A FLAVOR ✗');
  process.exit(allOk ? 0 : 1);
})();
