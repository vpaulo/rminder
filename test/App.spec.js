import { App } from '../src/App.js';

describe('App', () => {
  afterEach(() => {
    td.reset();
  });

  it('Should return App description', () => {
    expect(App()).equal('App initialiser');
  });

  it('Should post message', () => {
    td.replace(App.db, 'postMessage');

    App.start();

    td.verify(App.db.postMessage({ type: 'start' }), { times: 1 });
  });

  it('Should call dbMessage', () => {
    td.replace(App.rminder, 'tasks');
    const event = {
      data: {
        type: 'tasks',
      },
    };

    App.db.onmessage(event);

    td.verify(App.rminder.tasks(), { ignoreExtraArgs: true, times: 1 });
  });

  it('Should call screenTest', () => {
    td.replace(App.rminder, 'screenTest');

    App.start();

    td.verify(App.rminder.screenTest(), { times: 1 });
  });

  it('Should not remove expanded class from sidebar', () => {
    App.rminder.sidebar = {
      classList: {
        remove: td.func(),
      },
    };

    App.start();

    td.verify(App.rminder.sidebar.classList.remove('expanded'), { times: 0 });
  });

  it('Should remove expanded class from sidebar', () => {
    App.rminder.smallMediaQuery = {
      matches: true,
    };
    App.rminder.sidebar = {
      classList: {
        remove: td.func(),
      },
    };

    App.start();

    td.verify(App.rminder.sidebar.classList.remove('expanded'), { times: 1 });
  });

  it('Should call setDocHeight', () => {
    td.replace(App.rminder, 'setDocHeight');

    App.start();

    td.verify(App.rminder.setDocHeight(), { times: 1 });
  });
});
