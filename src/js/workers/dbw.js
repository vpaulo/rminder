const rminder = {
  name: 'rminder',
  version: 2, // Use a long long for this value (don't use a float)
  storeNames: ['tasks', 'settings'],
  settings: { completed: 'show', list: 'tasks', filter: 'oldest' },
};

const titles = {
  my_day: 'My day',
  important: 'Important',
  completed: 'Completed',
};

function openDb() {
  console.log('openDb ...');
  const req = indexedDB.open(rminder.name, rminder.version);
  req.onsuccess = () => {
    rminder.db = req.result;
    console.log('openDb DONE');
    useDB();
    postMessage({ type: 'opened', message: 'DB opened' });

    getSettings();
  };
  req.onerror = (evt) => {
    console.error('openDb:', evt.target.error);
  };

  req.onblocked = (evt) => {
    console.error(
      'openDb: Please close all other tabs with the App open',
      evt.target.error
    );
  };

  req.onupgradeneeded = () => {
    console.log('openDb.onupgradeneeded');
    rminder.db = req.result;
    if (!rminder.db.objectStoreNames.contains(rminder.storeNames[0])) {
      const store = rminder.db.createObjectStore(rminder.storeNames[0], {
        keyPath: 'id',
        autoIncrement: true,
      });

      store.createIndex('title', 'title', { unique: false });
      store.createIndex('important', 'important', { unique: false });
      store.createIndex('my_day', 'my_day', { unique: false });
      store.createIndex('completed', 'completed', { unique: false });
      store.createIndex('note', 'note', { unique: false });
      store.createIndex('creation_date', 'creation_date', { unique: false });
    }

    if (!rminder.db.objectStoreNames.contains(rminder.storeNames[1])) {
      const settings = rminder.db.createObjectStore(rminder.storeNames[1], {
        keyPath: 'id',
        autoIncrement: true,
      });

      settings.createIndex('completed', 'completed', { unique: false });
      settings.createIndex('list', 'list', { unique: false });
      settings.createIndex('filter', 'filter', { unique: false });
    }

    useDB();
  };
}

function closeDb() {
  rminder.db.close();
}

function useDB() {
  rminder.db.onversionchange = (evt) => {
    rminder.db.close();
    console.log(
      'openDb: A new version of this page is ready. Please reload or close this tab!',
      evt
    );
  };
}

/**
 * @param {string} store_name
 * @param {string} mode either "readonly" or "readwrite"
 */
function getObjectStore(store_name, mode) {
  return rminder.db.transaction(store_name, mode).objectStore(store_name);
}

function addTask(title, creation_date) {
  const obj = { title, creation_date };
  const store = getObjectStore(rminder.storeNames[0], 'readwrite');
  let req;
  try {
    if (rminder.settings.list === 'my_day') {
      obj.my_day = true;
    }
    if (rminder.settings.list === 'important') {
      obj.important = true;
    }
    req = store.add(obj);
  } catch (e) {
    if (e.name == 'DataCloneError') {
      postMessage({
        type: 'failure',
        // eslint-disable-next-line quotes
        message: "This engine doesn't know how to clone a Blob, use Firefox",
      });
    }

    throw e;
  }
  req.onsuccess = () => {
    postMessage({ type: 'success', message: 'addTask: successful' });
    displayTasks(store);
  };
  req.onerror = () => {
    postMessage({ type: 'failure', message: `addTask: error -> ${req.error}` });
  };
}

function removeTaskById(id) {
  const store = getObjectStore(rminder.storeNames[0], 'readwrite');
  store.openCursor().onsuccess = (evt) => {
    const cursor = evt.target.result;
    if (cursor) {
      if (cursor.value.id === id) {
        const request = cursor.delete();
        request.onsuccess = () => {
          postMessage({ type: 'success', message: `Task(${id}): deleted` });
          postMessage({ type: 'hideDetails' });
        };
      }
      // Move on to the next object in store
      cursor.continue();
    } else {
      displayTasks(store);
    }
  };
}

function showDetails(id) {
  const store = getObjectStore(rminder.storeNames[0], 'readonly');
  store.openCursor().onsuccess = (evt) => {
    const cursor = evt.target.result;
    if (cursor) {
      if (cursor.value.id === id) {
        postMessage({ type: 'details', key: cursor.key, value: cursor.value });
      }
      // Move on to the next object in store
      cursor.continue();
    }
  };
}

function updateTask(id, field, fieldValue) {
  const store = getObjectStore(rminder.storeNames[0], 'readwrite');
  store.openCursor().onsuccess = (evt) => {
    const cursor = evt.target.result;
    if (cursor) {
      const newValue = !fieldValue || cursor.value[field] !== fieldValue; // if field value undefined return true, otherwise compare values
      if (cursor.value.id === id && newValue) {
        const updateData = cursor.value;

        updateData[field] = fieldValue || !updateData[field]; // if field value undefined treat it as a boolean update

        const request = cursor.update(updateData);
        request.onsuccess = () => {
          postMessage({
            type: 'success',
            message: `Task(${id}): ${field} = ${updateData[field]}`,
          });
        };
      }
      // Move on to the next object in store
      cursor.continue();
    } else {
      displayTasks(store);
    }
  };
}

function displayTasks(store, settings = rminder.settings) {
  if (typeof store == 'undefined') {
    store = getObjectStore(rminder.storeNames[0], 'readonly');
  }

  // postMessage({ type: 'clear', message: 'Clear' });
  postMessage({ type: 'settings', settings });

  store.getAll().onsuccess = (e) => {
    const result =
      settings?.completed === 'hide'
        ? e.target.result.filter((task) => !task.completed)
        : e.target.result;
    const sortedResult = sortTasks(result, settings?.filter);
    let filteredList = sortedResult;
    if (settings.list !== 'tasks') {
      filteredList = filteredList.filter((task) => task[settings.list]);
      postMessage({
        type: 'tasks',
        value: sortedResult,
        list: {
          title: titles[settings.list],
          name: settings.list,
          value: filteredList,
        },
      });
    } else {
      postMessage({ type: 'tasks', value: sortedResult });
    }
  };
}

function updateSettings(completed, list, filter) {
  const cmpl = completed ? 'hide' : 'show';
  const settings = getObjectStore(rminder.storeNames[1], 'readwrite');

  settings.openCursor().onsuccess = (evt) => {
    const cursor = evt.target.result;
    if (cursor) {
      const updateSettings = cursor.value;

      if (completed !== undefined) {
        updateSettings.completed = cmpl;
      }

      if (list !== undefined) {
        updateSettings.list = list;
      }

      if (filter !== undefined) {
        updateSettings.filter = filter;
      }

      rminder.settings = updateSettings;

      const request = cursor.update(updateSettings);
      request.onsuccess = () => {
        postMessage({ type: 'success', message: 'Settings updated' });
      };
      // Move on to the next object in store
      cursor.continue();
    } else {
      displayTasks(undefined, rminder.settings);
    }
  };
}

function getSettings() {
  const settings = getObjectStore(rminder.storeNames[1], 'readonly');

  settings.getAll().onsuccess = (e) => {
    const result = e.target.result;
    if (e.target.result.length === 0) {
      setDefaultSettings();
    } else {
      rminder.settings = result[0];
      postMessage({ type: 'selectList', list: result[0].list });
      postMessage({ type: 'initialiseSettings', settings: result[0] });
    }
  };
}

function setDefaultSettings() {
  const settings = getObjectStore(rminder.storeNames[1], 'readwrite');
  const req = settings.add(rminder.settings);
  req.onsuccess = () => {
    postMessage({
      type: 'success',
      message: 'Set default settings: successful',
    });
    postMessage({ type: 'initialiseSettings', settings: rminder.settings });
  };
  req.onerror = () => {
    postMessage({
      type: 'failure',
      message: `Set default settings: error -> ${req.error}`,
    });
  };
}

function sortTasks(tasks, filter) {
  if (!filter) {
    return tasks;
  }

  let tsk;

  switch (filter) {
    case 'important':
      tsk = tasks.sort((x, y) =>
        x.important === y.important ? 0 : x.important ? -1 : 1
      ); // return important tasks first, for important tasks last (x.important === y.important) ? 0 : x.important ? 1 : -1
      break;
    case 'newest':
      tsk = tasks.sort((x, y) => y.creation_date - x.creation_date); // return newest tasks first
      break;
    default:
      tsk = tasks;
      break;
  }
  return tsk;
}

onmessage = (e) => {
  const { type } = e.data;

  switch (type) {
    case 'start':
      openDb();
      break;
    case 'close':
      closeDb();
      break;
    // case 'getBlob':
    // 		getBlob(e.data.key);
    // 		break;
    case 'addTask':
      addTask(e.data.title, e.data.creationDate);
      break;
    case 'removeTask':
      removeTaskById(e.data.id);
      break;
    case 'renameTask':
      updateTask(e.data.id, 'title', e.data.title);
      break;
    case 'showDetails':
      showDetails(e.data.id);
      break;
    case 'importantTask':
      updateTask(e.data.id, 'important');
      break;
    case 'myDayTask':
      updateTask(e.data.id, 'my_day');
      break;
    case 'noteTask':
      updateTask(e.data.id, 'note', e.data.note);
      break;
    case 'completedTask':
      updateTask(e.data.id, 'completed');
      break;
    // case 'clear':
    // 		clearObjectStore();
    // 		break;
    case 'settings':
      updateSettings(e.data.completed);
      break;
    case 'list':
      updateSettings(undefined, e.data.list);
      break;
    case 'filter':
      updateSettings(undefined, undefined, e.data.filter);
      break;
    case 'display':
      displayTasks();
      break;
    default:
      postMessage({ type });
      break;
  }
};
