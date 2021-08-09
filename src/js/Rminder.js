import { logger } from './logger.js';

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
		this.importanceCheckBtn = this.detailsContainer?.querySelector('.importance-check');
		this.completedCheck = this.detailsContainer?.querySelector('.completed-ckeck');
		this.note = document.querySelector('.note');
		this.listMyDay = document.querySelector('vp-nav-list[name="my_day"]');
		this.listImportant = document.querySelector('vp-nav-list[name="important"]');
		this.listCompleted = document.querySelector('vp-nav-list[name="completed"]');
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

	success(data) {
		const msg = typeof data.message !== 'undefined' ? `Success: ${data.message}` : 'Success: ';
		logger(msg, data);
	}

	opened(data, db) {
		logger(data.message);
		this.addEventListeners(db);
	}

	clear(data) {
		logger(data.message);
		this.taskInput.value = '';
		this.taskList.innerHTML = '';
	}

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

		// Show tasks
		this.taskList.innerHTML = dt.map(task => `<vp-list-task title="${task.title}" task="${task.id}" complete="${task.completed}" important="${task.important}" order="0"></vp-list-task>`).join('');
		// Show counters
		this.listMyDay.setAttribute('count', data.filter(d => d.my_day).length);
		this.listImportant.setAttribute('count', data.filter(d => d.important).length);
		this.listCompleted.setAttribute('count', data.filter(d => d.completed).length);
		this.listTasks.setAttribute('count', data.length); // TODO: change total for tasks list after adding lists functionality
		// Set List title
		this.listTitle.innerText = title;
		this.mainContainer.dataset.list = name;

		this.setDetailClasses(data);
	}

	details({ value: data }) {
		logger('Details: ', data);
		const date = new Date(data.creation_date);
		this.detailsContainer.setAttribute('aria-label', `Detail for task: ${data.title}`);
		this.detailsContainer.dataset.id = data.id;
		this.titleInput.value = data.title;
		this.note.value = data.note || '';
		this.creationDate.innerText = new Intl.DateTimeFormat().format(date);

		this.setDetailClasses(data);
	}

	addEventListeners(db) {
		this.addTaskBtn.addEventListener('click', () => { this.addTask(db); }, false);

		// TODO: create common function to handle completed check events
		this.completedCheck.addEventListener('click', e => {
			const parent = e.target.closest('[data-id]');
			if (this.toggleCompleted.checked) {
				this.hideDetails();
			}
			this.handleEvent('completedTask', db, +parent.dataset.id);
		}, false);

		this.lists.addEventListener('vp-nav-list:click', e => {
			const list = e.detail.name;
			if (list) {
				this.showList(list, db);

				if (this.smallMediaQuery.matches) {
					this.hideSidebar();
				}
			}
		});

		this.taskInput.addEventListener('keyup', event => {
			if (event.code === 'Enter') {
				this.addTask(db);
			}
		}, false);

		this.titleInput.addEventListener('keyup', event => {
			if (event.code === 'Enter') {
				this.renameTask(db);
			}
		}, false);

		// TODO: move the event listener to a common parent element
		document.addEventListener('vp-button:click', (e) => {
			if(e.detail.trigger) {
				this[e.detail.trigger](e.target, db);
			}
		});

		this.importanceCheckBtn.addEventListener('click', (e) => {
			const parent = e.target.closest('[data-id]');
			this.handleEvent('importantTask', db, +parent.dataset.id);
		}, false);

		this.myDay.addEventListener('click', () => {
			this.handleEvent('myDayTask', db);
		}, false);

		this.mediaQueryList.addEventListener('change', this.screenTest.bind(this), false);

		this.toggleCompleted.addEventListener('click', (evt) => { this.settingsCompleted(evt.target, db); }, false);

		this.settingsBtn.addEventListener('click', () => this.toggle(this.settingsBtn), false);
		this.filterBtn.addEventListener('click', () => this.toggle(this.filterBtn), false);

		this.orderFilters.forEach(filter => filter.addEventListener('change', evt => { this.filterUpdate(evt.target, db); }, false));

		window.addEventListener('resize', this.setDocHeight, false);
		window.addEventListener('orientationchange', this.setDocHeight, false);
	}

	addTask(db) {
		const title = this.taskInput.value.trim();
		const creationDate = Date.now();
		const list = this.mainContainer.dataset.list;
		if (title) {
			db.postMessage({ type: 'addTask', title, creationDate, list });
		} else {
			logger('Required field(s) missing: title');
		}
	}

	renameTask(elem, db) {
		const title = this.titleInput.value.trim();
		const id = +this.detailsContainer.dataset.id; // convert id to number
		const list = this.mainContainer.dataset.list;
		if (title) {
			db.postMessage({ type: 'renameTask', id, title, list });
		} else {
			logger('Required field(s) missing: title');
		}
	}

	completedTask(elem, db) {
			if (this.toggleCompleted.checked) {
				this.hideDetails();
			}
			this.handleEvent('completedTask', db, +elem.task);
	}

	importantTask(elem, db) {
		this.handleEvent('importantTask', db, +elem.task);
	}

	removeTask(elem, db) {
		this.handleEvent('removeTask', db);
	}

	showDetails(elem, db) {
		const id = +elem.task;
		if (id !== +this.detailsContainer.dataset.id) {
			db.postMessage({ type: 'showDetails', id });
		}
		this.detailsContainer.classList.add('expanded');
		this.screenTest();
		this.setSelected(elem);
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

	handleEvent(type, db, id = +this.detailsContainer.dataset.id) {
		const list = this.mainContainer.dataset.list;
		db.postMessage({ type, id, list });
	}

	setTaskNote(elem, db) {
		const id = +this.detailsContainer.dataset.id;
		const text = this.note.value.trim();
		const list = this.mainContainer.dataset.list;
		db.postMessage({ type: 'noteTask', id, note: text, list });
	}

	selectList({ list }) {
		document.querySelector('vp-nav-list.selected')?.classList?.remove('selected');
		document.querySelector(`vp-nav-list[name="${list}"]`).classList.add('selected');
	}

	showList(list, db) {
		this.selectList({ list });
		db.postMessage({ type: 'list', list });
	}

	setDetailClasses(data) {
		const id = +this.detailsContainer.dataset.id;
		let dt = data;

		if (Array.isArray(data)) {
			dt = data.find(d => d.id === id);
		}

		this.detailsContainer.classList.remove('important', 'completed', 'today');

		this.importanceCheckBtn?.querySelector('vp-icon')?.setAttribute('icon', `icons:star${ dt?.important ? '-solid' : ''}`);
		this.completedCheck?.querySelector('vp-icon')?.setAttribute('icon', `icons:${ dt?.completed ? 'check-' : ''}square`);

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

	screenTest(mql = this.mediaQueryList, elem = this.sidebar) {
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
		document.documentElement.style.setProperty('--vh', `${window.innerHeight}px`);
	}

	setSelected(el) {
		this.taskList.querySelector('.selected')?.classList?.remove('selected');
		el.classList.add('selected');
	}

	settingsCompleted(elem, db) {
		const checked = elem.checked;
		const list = this.mainContainer.dataset.list;
		db.postMessage({ type: 'settings', completed: checked, list });
	}

	toggle(elem) {
		elem.classList.toggle('open');
	}

	settings({ settings } = {}) {
		if (!settings) {
			return false;
		}

		if (settings.completed === 'hide') {
			this.toggleCompleted.checked = true;
			this.listCompleted.classList.add('hidden');
			if (this.listCompleted.classList.contains('selected')) {
				this.listTasks.click(); // Select Tasks list if Completed list was selected before hidding
			}
		} else {
			this.toggleCompleted.checked = false;
			this.listCompleted.classList.remove('hidden');
		}

		if (settings.filter) {
			this.orderFilters.forEach(filter => {
				if (filter.value === settings.filter) {
					filter.parentNode.click();
				}
			});
		}
	}

	filterUpdate(elem, db) {
		db.postMessage({ type: 'filter', filter: elem.value });
	}

	openModal() {
		this.modal.classList.add('open');
	}

	closeModal() {
		this.modal.classList.remove('open');
	}
}
