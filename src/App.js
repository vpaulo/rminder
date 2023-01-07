import { Rminder } from './js/Rminder.js';
import { dbMessage } from './js/dbMessage.js';

export const App = () => 'App initialiser';
App.rminder = new Rminder();
App.db = new Worker('./js/workers/dbw.js');
App.start = () => {
  App.db.postMessage({ type: 'start' });
  App.db.onmessage = (event) => dbMessage(App.rminder, event.data, App.db);

  if (App.rminder.smallMediaQuery.matches) {
    App.rminder.sidebar.classList.remove('expanded');
  }

  App.rminder.screenTest();
  App.rminder.setDocHeight();
};

App.start();
