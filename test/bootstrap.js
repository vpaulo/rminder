import puppeteer from 'puppeteer';
import { expect } from 'chai';

const open = process.env.OPEN;

const globalVariables = {
  browser: global.browser,
  expect: global.expect,
};

// puppeteer options
const opts = {
  headless: open ? false : true,
  slowMo: open ? 300 : 0,
  timeout: 5000,
  args: ['--start-maximized', '--window-size=1920,1040'],
};

// expose variables
before(async () => {
  global.expect = expect;
  global.browser = await puppeteer.launch(opts);
});

// close browser and reset global variables
after(async () => {
  await browser.close();

  global.browser = globalVariables.browser;
  global.expect = globalVariables.expect;
});

// Mock requests
// await page.setRequestInterception(true);
// page.on('request', async request => {

//   if (request.url().endsWith("/api/backendcall")) {
//     await request.respond({
//       status: 200,
//       contentType: 'application/json',
//       body: JSON.stringify({ key: 'value' })
//     });
//   } else {
//     await request.continue();
//   }
// });
