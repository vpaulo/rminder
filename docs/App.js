import { Rminder as p } from './js/Rminder.js';
import { dbMessage as e } from './js/dbMessage.js';
export const App = () => 'App initialiser';
App.rminder = new p();
App.db = new Worker('./js/workers/dbw.js');
App.start = () => {
  App.db.postMessage({ type: 'start' });
  App.db.onmessage = (p) => e(App.rminder, p.data, App.db);
  if (App.rminder.smallMediaQuery.matches) {
    App.rminder.sidebar.classList.remove('expanded');
  }
  App.rminder.screenTest();
  App.rminder.setDocHeight();
};
App.start();
