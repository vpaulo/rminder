import { logger as t } from './logger.js';
export class Rminder {
  constructor() {
    this.mediaQueryList = matchMedia('only screen and (max-width: 900px)');
    this.smallMediaQuery = matchMedia('only screen and (max-width: 630px)');
    this.taskInput = document.getElementById('task');
    this.addTaskBtn = document.querySelector('.add-task');
    this.taskList = document.querySelector('.tasks');
    this.detailsContainer = document.querySelector('.details');
    this.titleInput = document.querySelector('.title');
    this.myDay = document.querySelector('.my-day');
    this.importanceCheckBtn =
      this.detailsContainer?.querySelector('.importance-check');
    this.completedCheck =
      this.detailsContainer?.querySelector('.completed-ckeck');
    this.note = document.querySelector('.note');
    this.listMyDay = document.querySelector('vp-nav-list[name="my_day"]');
    this.listImportant = document.querySelector(
      'vp-nav-list[name="important"]'
    );
    this.listCompleted = document.querySelector(
      'vp-nav-list[name="completed"]'
    );
    this.listTasks = document.querySelector('vp-nav-list[name="tasks"]');
    this.sidebar = document.querySelector('.sidebar');
    this.lists = document.querySelector('.lists');
    this.listTitle = document.querySelector('.list-title');
    this.mainContainer = document.querySelector('.main');
    this.creationDate = document.querySelector('.creation-date');
    this.toggleCompleted = document.querySelector('.toggle-completed');
    this.settingsBtn = document.querySelector('.app__settings');
    this.filterBtn = document.querySelector('.list-filter');
    [...this.orderFilters] = document.querySelectorAll('.order-filter');
    this.modal = document.querySelector('vp-modal');
  }
  success(e) {
    const s =
      typeof e.message !== 'undefined' ? `Success: ${e.message}` : 'Success: ';
    t(s, e);
  }
  opened(e, s) {
    t(e.message);
    this.addEventListeners(s);
  }
  clear(e) {
    t(e.message);
    this.taskInput.value = '';
    this.taskList.innerHTML = '';
  }
  tasks({ value: e, list: s }) {
    t('Tasks:', e);
    let i = e;
    let a = 'Tasks';
    let n = 'tasks';
    if (s?.title) {
      i = s.value;
      a = s.title;
      n = s.name;
    }
    this.taskList.innerHTML = i
      .map(
        (t) =>
          `<vp-list-task title="${t.title}" task="${t.id}" complete="${t.completed}" important="${t.important}" order="0"></vp-list-task>`
      )
      .join('');
    this.listMyDay.setAttribute('count', e.filter((t) => t.my_day).length);
    this.listImportant.setAttribute(
      'count',
      e.filter((t) => t.important).length
    );
    this.listCompleted.setAttribute(
      'count',
      e.filter((t) => t.completed).length
    );
    this.listTasks.setAttribute('count', e.length);
    this.listTitle.innerText = a;
    this.mainContainer.dataset.list = n;
    this.setDetailClasses(e);
  }
  details({ value: e }) {
    t('Details: ', e);
    const s = new Date(e.creation_date);
    this.detailsContainer.setAttribute(
      'aria-label',
      `Detail for task: ${e.title}`
    );
    this.detailsContainer.dataset.id = e.id;
    this.titleInput.value = e.title;
    this.note.value = e.note || '';
    this.creationDate.innerText = new Intl.DateTimeFormat().format(s);
    this.setDetailClasses(e);
  }
  addEventListeners(t) {
    this.addTaskBtn.addEventListener(
      'click',
      () => {
        this.addTask(t);
      },
      false
    );
    this.completedCheck.addEventListener(
      'click',
      (e) => {
        const s = e.target.closest('[data-id]');
        if (this.toggleCompleted.checked) {
          this.hideDetails();
        }
        this.handleEvent('completedTask', t, +s.dataset.id);
      },
      false
    );
    this.lists.addEventListener('vp-nav-list:click', (e) => {
      const s = e.detail.name;
      if (s) {
        this.showList(s, t);
        if (this.smallMediaQuery.matches) {
          this.hideSidebar();
        }
      }
    });
    this.taskInput.addEventListener(
      'keyup',
      (e) => {
        if (e.code === 'Enter') {
          this.addTask(t);
        }
      },
      false
    );
    this.titleInput.addEventListener(
      'keyup',
      (e) => {
        if (e.code === 'Enter') {
          this.renameTask(t);
        }
      },
      false
    );
    document.addEventListener('vp-button:click', (e) => {
      if (e.detail.trigger) {
        this[e.detail.trigger](e.target, t);
      }
    });
    this.importanceCheckBtn.addEventListener(
      'click',
      (e) => {
        const s = e.target.closest('[data-id]');
        this.handleEvent('importantTask', t, +s.dataset.id);
      },
      false
    );
    this.myDay.addEventListener(
      'click',
      () => {
        this.handleEvent('myDayTask', t);
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
      (e) => {
        this.settingsCompleted(e.target, t);
      },
      false
    );
    this.settingsBtn.addEventListener(
      'click',
      () => this.toggle(this.settingsBtn),
      false
    );
    this.filterBtn.addEventListener(
      'click',
      () => this.toggle(this.filterBtn),
      false
    );
    this.orderFilters.forEach((e) =>
      e.addEventListener(
        'change',
        (e) => {
          this.filterUpdate(e.target, t);
        },
        false
      )
    );
    window.addEventListener('resize', this.setDocHeight, false);
    window.addEventListener('orientationchange', this.setDocHeight, false);
  }
  addTask(e) {
    const s = this.taskInput.value.trim();
    const i = Date.now();
    const a = this.mainContainer.dataset.list;
    if (s) {
      e.postMessage({ type: 'addTask', title: s, creationDate: i, list: a });
    } else {
      t('Required field(s) missing: title');
    }
  }
  renameTask(e, s) {
    const i = this.titleInput.value.trim();
    const a = +this.detailsContainer.dataset.id;
    const n = this.mainContainer.dataset.list;
    if (i) {
      s.postMessage({ type: 'renameTask', id: a, title: i, list: n });
    } else {
      t('Required field(s) missing: title');
    }
  }
  completedTask(t, e) {
    if (this.toggleCompleted.checked) {
      this.hideDetails();
    }
    this.handleEvent('completedTask', e, +t.task);
  }
  importantTask(t, e) {
    this.handleEvent('importantTask', e, +t.task);
  }
  removeTask(t, e) {
    this.handleEvent('removeTask', e);
  }
  showDetails(t, e) {
    const s = +t.task;
    if (s !== +this.detailsContainer.dataset.id) {
      e.postMessage({ type: 'showDetails', id: s });
    }
    this.detailsContainer.classList.add('expanded');
    this.screenTest();
    this.setSelected(t);
  }
  hideDetails() {
    this.detailsContainer.classList.remove('expanded');
    this.taskList.querySelector('.selected')?.classList?.remove('selected');
    this.modal.classList.remove('open');
    this.screenTest();
  }
  toggleSidebar() {
    this.sidebar.classList.toggle('expanded');
    this.screenTest(undefined, this.detailsContainer);
  }
  hideSidebar() {
    this.sidebar.classList.remove('expanded');
    this.screenTest();
  }
  handleEvent(t, e, s = +this.detailsContainer.dataset.id) {
    const i = this.mainContainer.dataset.list;
    e.postMessage({ type: t, id: s, list: i });
  }
  setTaskNote(t, e) {
    const s = +this.detailsContainer.dataset.id;
    const i = this.note.value.trim();
    const a = this.mainContainer.dataset.list;
    e.postMessage({ type: 'noteTask', id: s, note: i, list: a });
  }
  selectList({ list: t }) {
    document
      .querySelector('vp-nav-list.selected')
      ?.classList?.remove('selected');
    document
      .querySelector(`vp-nav-list[name="${t}"]`)
      .classList.add('selected');
  }
  showList(t, e) {
    this.selectList({ list: t });
    e.postMessage({ type: 'list', list: t });
  }
  setDetailClasses(t) {
    const e = +this.detailsContainer.dataset.id;
    let s = t;
    if (Array.isArray(t)) {
      s = t.find((t) => t.id === e);
    }
    this.detailsContainer.classList.remove('important', 'completed', 'today');
    this.importanceCheckBtn
      ?.querySelector('vp-icon')
      ?.setAttribute('icon', `icons:star${s?.important ? '-solid' : ''}`);
    this.completedCheck
      ?.querySelector('vp-icon')
      ?.setAttribute('icon', `icons:${s?.completed ? 'check-' : ''}square`);
    if (s?.important) {
      this.detailsContainer.classList.add('important');
    }
    if (s?.completed) {
      this.detailsContainer.classList.add('completed');
    }
    if (s?.my_day) {
      this.detailsContainer.classList.add('today');
    }
  }
  screenTest(t = this.mediaQueryList, e = this.sidebar) {
    if (t.matches && document.querySelectorAll('.expanded').length > 1) {
      e.classList.remove('expanded');
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
  setSelected(t) {
    this.taskList.querySelector('.selected')?.classList?.remove('selected');
    t.classList.add('selected');
  }
  settingsCompleted(t, e) {
    const s = t.checked;
    const i = this.mainContainer.dataset.list;
    e.postMessage({ type: 'settings', completed: s, list: i });
  }
  toggle(t) {
    t.classList.toggle('open');
  }
  settings({ settings: t } = {}) {
    if (!t) {
      return false;
    }
    if (t.completed === 'hide') {
      this.toggleCompleted.checked = true;
      this.listCompleted.classList.add('hidden');
      if (this.listCompleted.classList.contains('selected')) {
        this.listTasks.click();
      }
    } else {
      this.toggleCompleted.checked = false;
      this.listCompleted.classList.remove('hidden');
    }
    if (t.filter) {
      this.orderFilters.forEach((e) => {
        if (e.value === t.filter) {
          e.parentNode.click();
        }
      });
    }
  }
  filterUpdate(t, e) {
    e.postMessage({ type: 'filter', filter: t.value });
  }
  openModal() {
    this.modal.classList.add('open');
  }
  closeModal() {
    this.modal.classList.remove('open');
  }
}
