describe('Rminder App', () => {
  let page;

  afterEach(async () => {
    await page.close();
  });

  describe('Rminder - desktop', () => {
    beforeEach(async () => {
      page = await browser.newPage();
      await page.setViewport({ width: 1920, height: 1040 });
      await page.goto('http://127.0.0.1:5501/index.html', {
        waitUntil: 'networkidle2',
      });
    });

    it('should have the correct page title', async () => {
      expect(await page.title()).equals('Rminder');
    });

    it('should have sidebar open', async () => {
      const expanded = await page.$eval('.sidebar', (el) =>
        el.classList.contains('expanded')
      );
      expect(expanded).to.be.true;
    });

    it('should have details closed', async () => {
      const expanded = await page.$eval('.details', (el) =>
        el.classList.contains('expanded')
      );
      expect(expanded).to.be.false;
    });

    it('should not have tasks', async () => {
      const tasksCount = await page.$$eval(
        '.tasks > vp-list-task',
        (els) => els.length
      );
      expect(tasksCount).equals(0);
    });
  });

  describe('Rminder - mobile', () => {
    beforeEach(async () => {
      page = await browser.newPage();
      await page.setViewport({ width: 600, height: 800 });
      await page.goto('http://127.0.0.1:5501/index.html', {
        waitUntil: 'networkidle2',
      });
    });

    it('should have sidebar closed', async () => {
      const expanded = await page.$eval('.sidebar', (el) =>
        el.classList.contains('expanded')
      );
      expect(expanded).to.be.false;
    });
  });
});
