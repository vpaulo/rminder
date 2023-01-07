const e = {
  name: 'rminder',
  version: 2,
  storeNames: ['tasks', 'settings'],
  settings: { completed: 'show', list: 'tasks', filter: 'oldest' },
};
const t = { my_day: 'My day', important: 'Important', completed: 'Completed' };
function s() {
  console.log('openDb ...');
  const t = indexedDB.open(e.name, e.version);
  t.onsuccess = () => {
    e.db = t.result;
    console.log('openDb DONE');
    n();
    postMessage({ type: 'opened', message: 'DB opened' });
    p();
  };
  t.onerror = (e) => {
    console.error('openDb:', e.target.error);
  };
  t.onblocked = (e) => {
    console.error(
      'openDb: Please close all other tabs with the App open',
      e.target.error
    );
  };
  t.onupgradeneeded = () => {
    console.log('openDb.onupgradeneeded');
    e.db = t.result;
    if (!e.db.objectStoreNames.contains(e.storeNames[0])) {
      const t = e.db.createObjectStore(e.storeNames[0], {
        keyPath: 'id',
        autoIncrement: true,
      });
      t.createIndex('title', 'title', { unique: false });
      t.createIndex('important', 'important', { unique: false });
      t.createIndex('my_day', 'my_day', { unique: false });
      t.createIndex('completed', 'completed', { unique: false });
      t.createIndex('note', 'note', { unique: false });
      t.createIndex('creation_date', 'creation_date', { unique: false });
    }
    if (!e.db.objectStoreNames.contains(e.storeNames[1])) {
      const t = e.db.createObjectStore(e.storeNames[1], {
        keyPath: 'id',
        autoIncrement: true,
      });
      t.createIndex('completed', 'completed', { unique: false });
      t.createIndex('list', 'list', { unique: false });
      t.createIndex('filter', 'filter', { unique: false });
    }
    n();
  };
}
function a() {
  e.db.close();
}
function n() {
  e.db.onversionchange = (t) => {
    e.db.close();
    console.log(
      'openDb: A new version of this page is ready. Please reload or close this tab!',
      t
    );
  };
}
function o(t, s) {
  return e.db.transaction(t, s).objectStore(t);
}
function r(t, s, a) {
  const n = { title: t, creation_date: s };
  const r = o(e.storeNames[0], 'readwrite');
  let i;
  try {
    if (a === 'my_day') {
      n.my_day = true;
    }
    if (a === 'important') {
      n.important = true;
    }
    i = r.add(n);
  } catch (e) {
    if (e.name == 'DataCloneError') {
      postMessage({
        type: 'failure',
        message: "This engine doesn't know how to clone a Blob, use Firefox",
      });
    }
    throw e;
  }
  i.onsuccess = () => {
    postMessage({ type: 'success', message: 'addTask: successful' });
    l(r, a);
  };
  i.onerror = () => {
    postMessage({ type: 'failure', message: `addTask: error -> ${i.error}` });
  };
}
function i(t, s) {
  const a = o(e.storeNames[0], 'readwrite');
  a.openCursor().onsuccess = (e) => {
    const n = e.target.result;
    if (n) {
      if (n.value.id === t) {
        const e = n.delete();
        e.onsuccess = () => {
          postMessage({ type: 'success', message: `Task(${t}): deleted` });
          postMessage({ type: 'hideDetails' });
        };
      }
      n.continue();
    } else {
      l(a, s);
    }
  };
}
function c(t) {
  const s = o(e.storeNames[0], 'readonly');
  s.openCursor().onsuccess = (e) => {
    const s = e.target.result;
    if (s) {
      if (s.value.id === t) {
        postMessage({ type: 'details', key: s.key, value: s.value });
      }
      s.continue();
    }
  };
}
function d(t, s, a, n) {
  const r = o(e.storeNames[0], 'readwrite');
  r.openCursor().onsuccess = (e) => {
    const o = e.target.result;
    if (o) {
      const e = !a || o.value[s] !== a;
      if (o.value.id === t && e) {
        const e = o.value;
        e[s] = a || !e[s];
        const n = o.update(e);
        n.onsuccess = () => {
          postMessage({
            type: 'success',
            message: `Task(${t}): ${s} = ${e[s]}`,
          });
        };
      }
      o.continue();
    } else {
      l(r, n);
    }
  };
}
function l(s, a, n = e.settings) {
  if (typeof s == 'undefined') {
    s = o(e.storeNames[0], 'readonly');
  }
  postMessage({ type: 'clear', message: 'Clear' });
  postMessage({ type: 'settings', settings: n });
  s.getAll().onsuccess = (e) => {
    const s =
      n?.completed === 'hide'
        ? e.target.result.filter((e) => !e.completed)
        : e.target.result;
    const o = m(s, n?.filter);
    let r = o;
    if (a !== 'tasks') {
      r = r.filter((e) => e[a]);
      postMessage({
        type: 'tasks',
        value: o,
        list: { title: t[a], name: a, value: r },
      });
    } else {
      postMessage({ type: 'tasks', value: o });
    }
  };
}
function u(t, s, a) {
  const n = t ? 'hide' : 'show';
  const r = o(e.storeNames[1], 'readwrite');
  r.openCursor().onsuccess = (o) => {
    const r = o.target.result;
    if (r) {
      const o = r.value;
      if (t !== undefined) {
        o.completed = n;
      }
      if (s !== undefined) {
        o.list = s;
      }
      if (a !== undefined) {
        o.filter = a;
      }
      e.settings = o;
      const i = r.update(o);
      i.onsuccess = () => {
        postMessage({ type: 'success', message: 'Settings updated' });
      };
      r.continue();
    } else {
      l(undefined, s, e.settings);
    }
  };
}
function p() {
  const t = o(e.storeNames[1], 'readonly');
  t.getAll().onsuccess = (t) => {
    const s = t.target.result;
    if (t.target.result.length === 0) {
      f();
    } else {
      e.settings = s[0];
      l(undefined, s[0].list, s[0]);
      postMessage({ type: 'selectList', list: s[0].list });
    }
  };
}
function f() {
  const t = o(e.storeNames[1], 'readwrite');
  const s = t.add(e.settings);
  s.onsuccess = () => {
    postMessage({
      type: 'success',
      message: 'Set default settings: successful',
    });
    l(undefined, e.settings.list, e.settings);
  };
  s.onerror = () => {
    postMessage({
      type: 'failure',
      message: `Set default settings: error -> ${s.error}`,
    });
  };
}
function m(e, t) {
  if (!t) {
    return e;
  }
  let s;
  switch (t) {
    case 'important':
      s = e.sort((e, t) =>
        e.important === t.important ? 0 : e.important ? -1 : 1
      );
      break;
    case 'newest':
      s = e.sort((e, t) => t.creation_date - e.creation_date);
      break;
    default:
      s = e;
      break;
  }
  return s;
}
onmessage = (e) => {
  const { type: t } = e.data;
  switch (t) {
    case 'start':
      s();
      break;
    case 'close':
      a();
      break;
    case 'addTask':
      r(e.data.title, e.data.creationDate, e.data.list);
      break;
    case 'removeTask':
      i(e.data.id, e.data.list);
      break;
    case 'renameTask':
      d(e.data.id, 'title', e.data.title, e.data.list);
      break;
    case 'showDetails':
      c(e.data.id);
      break;
    case 'importantTask':
      d(e.data.id, 'important', undefined, e.data.list);
      break;
    case 'myDayTask':
      d(e.data.id, 'my_day', undefined, e.data.list);
      break;
    case 'noteTask':
      d(e.data.id, 'note', e.data.note, e.data.list);
      break;
    case 'completedTask':
      d(e.data.id, 'completed', undefined, e.data.list);
      break;
    case 'settings':
      u(e.data.completed, e.data.list);
      break;
    case 'list':
      u(undefined, e.data.list);
      break;
    case 'filter':
      u(undefined, undefined, e.data.filter);
      break;
    case 'display':
      l(undefined, e.data.list);
      break;
    default:
      postMessage({ type: t });
      break;
  }
};
