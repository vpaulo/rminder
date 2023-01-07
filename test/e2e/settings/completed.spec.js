describe('Rminder - settings', () => {
  let page;

  beforeEach(async () => {
    page = await browser.newPage();
    await page.setViewport({ width: 1920, height: 1040 });
    await page.goto('http://127.0.0.1:5501/index.html', {
      waitUntil: 'networkidle2',
    });
  });

  afterEach(async () => {
    await page.close();
  });

  it('should set completed tasks to be hidden', async () => {
    await page.click('.app__settings');

    await page.waitForSelector('.switch');
    await page.click('.switch');

    const checked = await page.$eval('.toggle-completed', (el) => el.checked);
    expect(checked).to.be.true;
  });

  it('should set completed tasks to be visible', async () => {
    await page.click('.app__settings');

    await page.waitForSelector('.switch');
    await page.click('.switch');

    const checked = await page.$eval('.toggle-completed', (el) => el.checked);
    expect(checked).to.be.false;
  });
});
