describe.skip('DB worker', () => {
  const db = new Worker('/base/js/workers/dbw.js');
  db.onmessage = (result) => result;

  beforeEach(() => {
    db.postMessage({ type: 'start' });
  });

  afterEach(() => {
    // db.postMessage({ type: 'close' });
    td.reset();
  });

  describe.skip('Open DB', () => {
    it('Should open db', () => {
      td.replace(db, 'onmessage');
      // db.onmessage = (result) => {
      // 	expect(result.data.type).equals('opened');
      // 	expect(result.data.message).equals('DB opened');
      // 	done();
      // }
      td.verify(db.onmessage());
    });
  });

  describe('Add task', () => {
    it('Should add task', (done) => {
      // db.postMessage({ type: 'start' });
      db.postMessage({
        type: 'addTask',
        title: 'task one',
        creationDate: 1617899343187,
        list: '',
      });

      db.onmessage = (result) => {
        expect(result.data.type).equals('success');
        done();
      };
    });
  });
});
