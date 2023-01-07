import { dbMessage } from '../../src/js/dbMessage.js';

describe('dbMessage', () => {
  const foo = {
    bar: td.func(),
  };

  beforeEach(() => {
    td.replace(console, 'log');
  });

  afterEach(() => {
    td.reset();
    window.location.hash = '';
  });

  it('Should log message if method does not exist', () => {
    window.location.hash = 'debug';
    dbMessage(foo, { type: 'test' }, {});

    td.verify(
      console.log('Error running db worker - type: test does not exist', [])
    );
  });

  it('Should call method', () => {
    const data = { type: 'bar' };
    const db = {};

    dbMessage(foo, data, db);

    td.verify(foo.bar(data, db));
  });
});
