import { logger } from './logger.js';

export class Rminder {
  constructor() {
    this.mediaQueryList = matchMedia('only screen and (max-width: 900px)');
    this.smallMediaQuery = matchMedia('only screen and (max-width: 630px)');

    this.detailsContainer = document.querySelector('.details');
    this.titleInput = document.querySelector('.title');
    this.myDay = document.querySelector('.my-day');
    this.importanceCheckBtn =
      this.detailsContainer?.querySelector('.importance-check');
    this.completedCheck =
      this.detailsContainer?.querySelector('.completed-ckeck');
    this.note = document.querySelector('.note');
    this.drawer = document.querySelector('vp-rminder-tasks-drawer');
    this.creationDate = document.querySelector('.creation-date');
    this.toggleCompleted = document.querySelector('.toggle-completed');
    this.settingsBtn = document.querySelector('.app__settings');
    this.modal = document.querySelector('vp-modal');
    this.rminderList = document.querySelector('vp-rminder-list');
  }

  success(data) {
    const msg =
      typeof data.message !== 'undefined'
        ? `Success: ${data.message}`
        : 'Success: ';
    logger(msg, data);
  }

  opened(data, db) {
    logger(data.message);
    this.addEventListeners(db);
  }

  initialiseSettings(data) {
    this.rminderList.setAttribute('list', data.settings.list);
    this.rminderList.setAttribute('filter', data.settings.filter);
  }

  // clear(data) {
  // 	logger(data.message);
  // 	this.taskList.innerHTML = '';
  // }

  tasks({ value: data, list }) {
    logger('Tasks:', data);
    let dt = data;
    let title = 'Tasks';
    let name = 'tasks';

    if (list?.title) {
      dt = list.value;
      title = list.title;
      name = list.name;
    }

    // Show counters
    this.drawer.updateTasks(data);

    this.setDetailClasses(data);

    this.rminderList.dispatchEvent(
      new CustomEvent('vp:tasks', {
        detail: {
          tasks: dt,
          list: name,
          title,
        },
        bubbles: true,
        composed: true,
      })
    );
  }

  details({ value: data }) {
    logger('Details: ', data);
    const date = new Date(data.creation_date);
    this.detailsContainer.setAttribute(
      'aria-label',
      `Detail for task: ${data.title}`
    );
    this.detailsContainer.dataset.id = data.id;
    this.titleInput.value = data.title;
    this.note.value = data.note || '';
    this.creationDate.innerText = new Intl.DateTimeFormat().format(date);

    this.setDetailClasses(data);
  }

  addEventListeners(db) {
    // TODO: create common function to handle completed check events
    this.completedCheck.addEventListener(
      'click',
      (e) => {
        const parent = e.target.closest('[data-id]');
        if (this.toggleCompleted.checked) {
          this.hideDetails();
        }
        this.handleEvent('completedTask', db, +parent.dataset.id);
      },
      false
    );

    this.drawer.addEventListener('vp-nav-list:click', (e) => {
      const list = e.detail.name;
      if (list) {
        this.showList(list, db);

        if (this.smallMediaQuery.matches) {
          this.hideSidebar();
        }
      }
    });

    this.titleInput.addEventListener(
      'keyup',
      (event) => {
        if (event.code === 'Enter') {
          this.renameTask(db);
        }
      },
      false
    );

    // TODO: move the event listener to a common parent element
    // TODO: showDetails trigger throws error
    // TODO: it adds extra task from add-list-task
    document.addEventListener('vp:trigger', (e) => {
      if (e.detail.trigger) {
        this[e.detail.trigger](e.target, db, e.detail);
      }
    });

    document.addEventListener('modal', (e) => {
      this.eventHandler(e, db);
    });
    document.addEventListener('vp:addTask', (e) => {
      this.addTask(e.detail, db);
    });
    document.addEventListener('vp:showList', (e) => {
      this.showList(e.detail.list, db);
    });
    document.addEventListener('vp:updateTask', (e) => {
      this[e.detail.type](e.detail.id, db, e.detail);
    });
    document.addEventListener('vp:filterUpdate', (e) => {
      this.filterUpdate(e.detail.filter, db);
    });

    this.importanceCheckBtn.addEventListener(
      'click',
      (e) => {
        const parent = e.target.closest('[data-id]');
        this.handleEvent('importantTask', db, +parent.dataset.id);
      },
      false
    );

    this.myDay.addEventListener(
      'click',
      () => {
        this.handleEvent('myDayTask', db);
      },
      false
    );

    this.mediaQueryList.addEventListener(
      'change',
      this.screenTest.bind(this),
      false
    );

    this.toggleCompleted.addEventListener(
      'click',
      (evt) => {
        this.settingsCompleted(evt.target, db);
      },
      false
    );

    this.settingsBtn.addEventListener(
      'click',
      () => this.toggle(this.settingsBtn),
      false
    );

    window.addEventListener('resize', this.setDocHeight, false);
    window.addEventListener('orientationchange', this.setDocHeight, false);
  }

  addTask(data, db) {
    db.postMessage({
      type: 'addTask',
      title: data.title,
      creationDate: data.creationDate,
    });
  }

  renameTask(elem, db) {
    const title = this.titleInput.value.trim();
    const id = +this.detailsContainer.dataset.id; // convert id to number
    if (title) {
      db.postMessage({
        type: 'renameTask',
        id,
        title,
      });
    } else {
      logger('Required field(s) missing: title');
    }
  }

  completedTask(id, db) {
    if (this.toggleCompleted.checked) {
      this.hideDetails();
    }
    this.handleEvent('completedTask', db, id);
  }

  importantTask(id, db) {
    this.handleEvent('importantTask', db, id);
  }

  removeTask(elem, db) {
    this.handleEvent('removeTask', db);
  }

  showDetails(id, db) {
    if (id !== +this.detailsContainer.dataset.id) {
      db.postMessage({
        type: 'showDetails',
        id,
      });
    }
    this.detailsContainer.classList.add('expanded');
    this.screenTest();
  }

  hideDetails() {
    this.detailsContainer.classList.remove('expanded');
    // this.taskList.querySelector('.selected')?.classList?.remove('selected');
    this.modal.close();
    this.screenTest();
  }

  toggleSidebar() {
    // TODO: this.drawer.toggle(); should work
    this.drawer.classList.toggle('expanded');
    this.screenTest(undefined, this.detailsContainer);
  }

  hideSidebar() {
    // TODO: maybe create a Drawer.hide() method
    this.drawer.classList.remove('expanded');
    this.screenTest();
  }

  handleEvent(type, db, id = +this.detailsContainer.dataset.id) {
    db.postMessage({
      type,
      id,
    });
  }

  setTaskNote(elem, db) {
    const id = +this.detailsContainer.dataset.id;
    const text = this.note.value.trim();
    db.postMessage({
      type: 'noteTask',
      id,
      note: text,
    });
  }

  selectList({ list }) {
    document
      .querySelector('vp-nav-list.selected')
      ?.classList.remove('selected');
    document
      .querySelector(`vp-nav-list[name="${list}"]`)
      ?.classList.add('selected');
  }

  showList(list, db) {
    this.selectList({
      list,
    });
    db.postMessage({
      type: 'list',
      list,
    });
  }

  setDetailClasses(data) {
    const id = +this.detailsContainer.dataset.id;
    let dt = data;

    if (Array.isArray(data)) {
      dt = data.find((d) => d.id === id);
    }

    this.detailsContainer.classList.remove('important', 'completed', 'today');

    this.importanceCheckBtn
      ?.querySelector('vp-icon')
      ?.setAttribute('icon', `icons:star${dt?.important ? '-solid' : ''}`);
    this.completedCheck
      ?.querySelector('vp-icon')
      ?.setAttribute('icon', `icons:${dt?.completed ? 'check-' : ''}square`);

    // TODO: remove this classes
    if (dt?.important) {
      this.detailsContainer.classList.add('important');
    }

    if (dt?.completed) {
      this.detailsContainer.classList.add('completed');
    }

    if (dt?.my_day) {
      this.detailsContainer.classList.add('today');
    }
  }

  screenTest(mql = this.mediaQueryList, elem = this.drawer) {
    if (mql.matches && document.querySelectorAll('.expanded').length > 1) {
      elem.classList.remove('expanded');
    }

    if (this.smallMediaQuery.matches && document.querySelector('.expanded')) {
      this.mainContainer?.classList.add('hidden');
    } else {
      this.mainContainer?.classList.remove('hidden');
    }
  }

  setDocHeight() {
    document.documentElement.style.setProperty(
      '--vh',
      `${window.innerHeight}px`
    );
  }

  // setSelected(el) {
  // 	this.taskList.querySelector('.selected')?.classList?.remove('selected');
  // 	el.classList.add('selected');
  // }

  settingsCompleted(elem, db) {
    const checked = elem.checked;
    db.postMessage({
      type: 'settings',
      completed: checked,
    });
  }

  toggle(elem) {
    elem.classList.toggle('open');
  }

  // settings({ settings } = {}) {
  // 	if (!settings) {
  // 		return false;
  // 	}

  // 	if (settings.completed === 'hide') {
  // 		this.toggleCompleted.checked = true;
  // 		this.listCompleted.classList.add('hidden');
  // 		if (this.listCompleted.classList.contains('selected')) {
  // 			this.listTasks.click(); // Select Tasks list if Completed list was selected before hidding
  // 		}
  // 	} else {
  // 		this.toggleCompleted.checked = false;
  // 		this.listCompleted.classList.remove('hidden');
  // 	}

  // 	if (settings.filter) {
  // 		this.rminderList.setAttribute('filter', settings.filter);
  // 	}
  // }

  filterUpdate(filter, db) {
    db.postMessage({
      type: 'filter',
      filter,
    });
  }

  openModal() {
    this.modal.open();
  }

  closeModal() {
    this.modal.close();
  }

  eventHandler(event, db) {
    if (event.detail.trigger) {
      this[event.detail.trigger](event.target, db, event.detail);
    }
  }
}
