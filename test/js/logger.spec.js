import { logger } from '../../src/js/logger.js';

describe('logger', () => {
  beforeEach(() => {
    td.replace(console, 'log');
  });

  afterEach(() => {
    td.reset();
    window.location.hash = '';
  });

  it('Should not log if not in debug mode', () => {
    logger('test');
    td.verify(console.log('test', []), { times: 0 });
  });

  it('Should log if in debug mode', () => {
    window.location.hash = 'debug';
    logger('test');
    td.verify(console.log('test', []), { times: 1 });
  });

  it('Should log with message and args', () => {
    window.location.hash = 'debug';
    logger('test', 'one', 2);
    td.verify(console.log('test', ['one', 2]), { times: 1 });
  });
});
