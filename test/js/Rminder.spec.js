import { Rminder } from '../../src/js/Rminder.js';

describe('Rminder', () => {
	let rminder;

	beforeEach(() => {
		window.location.hash = 'debug';

		document.body.innerHTML = `
		<header>
			<div class="header__start">
				<span class="app__logo"></span>
				<span class="app__name">RMINDER</span>
			</div>
			<span class="app__settings">
				<div class="settings">
					<label class="switch">
						<span class="text">Hide Completed</span>
						<input class="toggle-completed" type="checkbox">
						<span class="slider"></span>
					</label>
				</div>
			</span>
		</header>
		<main id="app">
			<!-- left colunm -->
			<aside class="sidebar expanded" aria-label="Lists menu">
				<div class="sidebar__header">
					<span class="menu" aria-label="Toggle sidebar"></span>
				</div>
				<div class="sidebar__content">
					<div class="lists" aria-label="Lists" role="navigation">
						<vp-nav-list name="my_day" label="My day" count="0" icon="icons:today"></vp-nav-list>
						<vp-nav-list name="important" label="Important" count="0" icon="icons:star"></vp-nav-list>
						<vp-nav-list name="completed" label="Completed" count="0" icon="icons:check-square"></vp-nav-list>
						<vp-nav-list name="tasks" label="Tasks" count="0" icon="icons:tasks"></vp-nav-list>
					</div>
				</div>
				<div class="sidebar__add-list">
					<div class="add-list"></div>
					<div class="add-group-list"></div>
				</div>
			</aside>
			<!-- center column -->
			<div class="main" role="main" data-list="">
				<div class="list-toolbar">
					<span class="list-title">Tasks</span>
					<span class="list-filter">
						<div class="filters">
							<span class="legend">Order by:</span>
							<label class="switch"><span class="text">Important</span> <input class="order-filter" name="order-filter" type="radio" value="important"> <span class="radio"></span></label>
							<label class="switch"><span class="text">Oldest</span> <input class="order-filter" name="order-filter" value="oldest" type="radio" checked> <span class="radio"></span></label>
							<label class="switch"><span class="text">Newest</span> <input class="order-filter" name="order-filter" value="newest" type="radio"> <span class="radio"></span></label>
						</div>
					</span>
				</div>
				<div class="container">
					<div class="add-tasks">
						<label class="add-task--label" for="task"></label>
						<input type="text" name="task" id="task" maxlength="255" aria-label="Add a task" placeholder="Add a task" />
						<button class="add-task">Add</button>
					</div>
					<div class="tasks"></div>
				</div>
			</div>
			<!-- right colunm -->
			<aside class="details" aria-label="Detail for task: {{task-selected}}" data-id="">
				<div class="details__body">
					<div class="detail__title">
						<span class="completed-ckeck" title="Set it as complete"></span>
						<input class="title" type="text" value="" />
						<vp-button class="rename" trigger="renameTask">Rename</vp-button>
					</div>
					<div class="detail__my-day">
						<span class="my-day" title="Add to my day"></span>
						<span class="importance-check" title="Set it as important"></span>
					</div>
					<div class="detail__note">
						<textarea class="note" cols="30" rows="5" placeholder="Add notes"></textarea>
						<vp-button class="add-note" trigger="setTaskNote">Add</vp-button>
					</div>
				</div>
				<div class="details__footer">
					<vp-button class="close"></vp-button>
					<span class="creation-date"></span>
					<vp-button class="remove"></vp-button>
				</div>
			</aside>
		</main>
		<vp-modal>
			<span slot="body">Task will be permanent deleted, you won't be able to undo this action.</span>
			<div slot="footer">
				<vp-button class="default" type="default" trigger="closeModal">Cancel</vp-button>
				<vp-button class="warning" type="danger" trigger="removeTask">Delete task</vp-button>
			</div>
		</vp-modal>
		`;

		rminder = new Rminder();
	});

	afterEach(() => {
		td.reset();
		window.location.hash = '';
		document.body.innerHTML = '';
	});

	describe('Rminder.success', () => {
		beforeEach(() => {
			td.replace(console, 'log');
		});

		it('Should log success message', () => {
			const data = { type: 'test', message: 'ok' };

			rminder.success(data);

			td.verify(console.log('Success: ok', [data]));
		});

		it('Should log success when message is undefined', () => {
			const data = {};

			rminder.success(data);

			td.verify(console.log('Success: ', [data]));
		});
	});

	describe('Rminder.opened', () => {
		it('Should call addEventListeners', () => {
			td.replace(rminder, 'addEventListeners');

			const data = { type: 'test', message: 'ok' };
			const db = {};

			rminder.opened(data, db);

			td.verify(rminder.addEventListeners(db));
		});
	});

	describe('Rminder.clear', () => {
		it('Should clear elements', () => {
			rminder.taskInput.value = 'test';
			rminder.taskList.innerHTML = '<li>one</li>';

			rminder.clear({});

			expect(rminder.taskInput.value).equals('');
			expect(rminder.taskList.innerHTML).equals('');
		});
	});

	describe('Rminder.tasks', () => {
		it('Should add task', () => {
			td.replace(rminder, 'setDetailClasses');

			const data = {
				value: [{
					id: 1,
					completed: false,
					important: false,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(data);

			expect(rminder.taskList.innerHTML).equals('<vp-list-task title="test" task="1" complete="false" important="false" order="0"></vp-list-task>');
			expect(rminder.listMyDay.getAttribute('count')).equals('0');
			expect(rminder.listImportant.getAttribute('count')).equals('0');
			expect(rminder.listTasks.getAttribute('count')).equals('1');
			expect(rminder.listTitle.innerText).equals('Tasks');
			expect(rminder.mainContainer.dataset.list).equals('tasks');

			td.verify(rminder.setDetailClasses(data.value));
		});

		it('Should add completed task', () => {
			td.replace(rminder, 'setDetailClasses');

			const data = {
				value: [{
					id: 1,
					completed: true,
					important: false,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(data);

			expect(rminder.taskList.innerHTML).equals('<vp-list-task title="test" task="1" complete="true" important="false" order="0"></vp-list-task>');

			td.verify(rminder.setDetailClasses(data.value));
		});

		it('Should add important task', () => {
			td.replace(rminder, 'setDetailClasses');

			const data = {
				value: [{
					id: 1,
					completed: false,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(data);

			expect(rminder.taskList.innerHTML).equals('<vp-list-task title="test" task="1" complete="false" important="true" order="0"></vp-list-task>');
			expect(rminder.listImportant.getAttribute('count')).equals('1');

			td.verify(rminder.setDetailClasses(data.value));
		});

		it('Should add completed and important task', () => {
			td.replace(rminder, 'setDetailClasses');

			const data = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(data);

			expect(rminder.taskList.innerHTML).equals('<vp-list-task title="test" task="1" complete="true" important="true" order="0"></vp-list-task>');
			expect(rminder.listImportant.getAttribute('count')).equals('1');

			td.verify(rminder.setDetailClasses(data.value));
		});

		it('Should add task to specific list', () => {
			td.replace(rminder, 'setDetailClasses');

			const data = {
				value: [{
					id: 1,
					completed: false,
					important: false,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}],
				list: {
					title: 'Important',
					name: 'important',
					value: [{
						id: 2,
						completed: false,
						important: true,
						title: 'test2',
						my_day: false,
						note: '',
						creation_date: 1616785736468
					}]
				}
			};

			rminder.tasks(data);

			expect(rminder.taskList.innerHTML).equals('<vp-list-task title="test2" task="2" complete="false" important="true" order="0"></vp-list-task>');
			expect(rminder.listTitle.innerText).equals('Important');
			expect(rminder.mainContainer.dataset.list).equals('important');

			td.verify(rminder.setDetailClasses(data.value));
		});
	});

	describe('Rminder.details', () => {
		it('Should fill details info', () => {
			td.replace(rminder, 'setDetailClasses');

			const data = {
				value: {
					id: 1,
					completed: false,
					important: false,
					title: 'test',
					my_day: false,
					note: 'something to take note',
					creation_date: 1616785736469
				}
			};

			rminder.details(data);

			expect(rminder.detailsContainer.getAttribute('aria-label')).equals('Detail for task: test');
			expect(rminder.detailsContainer.dataset.id).equals('1');
			expect(rminder.titleInput.value).equals('test');
			expect(rminder.note.value).equals('something to take note');
			expect(rminder.creationDate.innerText).equals(new Intl.DateTimeFormat().format(new Date(1616785736469)));

			td.verify(rminder.setDetailClasses(data.value));
		});

		it('Should show empty note if it is undefined', () => {
			td.replace(rminder, 'setDetailClasses');

			const data = {
				value: {
					id: 1,
					completed: false,
					important: false,
					title: 'test',
					my_day: false,
					creation_date: 1616785736469
				}
			};

			rminder.details(data);

			expect(rminder.note.value).equals('');

			td.verify(rminder.setDetailClasses(data.value));
		});
	});

	describe('Rminder.addTask', () => {
		it('Should send postMessage to add task', () => {
			const db = {
				postMessage: td.func()
			};

			const data = {
				type: 'addTask',
				title: 'test',
				creationDate: td.matchers.anything(),
				list: 'tasks'
			};

			rminder.taskInput.value = 'test';
			rminder.mainContainer.dataset.list = 'tasks';

			rminder.addTask(db);

			td.verify(db.postMessage(data));
		});

		it('Should not send postMessage if title is missing', () => {
			td.replace(console, 'log');
			const db = {
				postMessage: td.func()
			};

			rminder.taskInput.value = '';
			rminder.mainContainer.dataset.list = 'tasks';

			rminder.addTask(db);

			td.verify(db.postMessage(), { times: 0, ignoreExtraArgs: true });
			td.verify(console.log('Required field(s) missing: title', []), { times: 1 });
		});
	});

	describe('Rminder.renameTask', () => {
		it('Should send postMessage to rename task', () => {
			const db = {
				postMessage: td.func()
			};

			const data = {
				type: 'renameTask',
				id: 2,
				title: 'test2',
				list: ''
			};

			rminder.titleInput.value = 'test2';
			rminder.detailsContainer.dataset.id = '2';

			rminder.renameTask(rminder.rename, db);

			td.verify(db.postMessage(data));
		});

		it('Should not send postMessage if title is missing', () => {
			td.replace(console, 'log');
			const db = {
				postMessage: td.func()
			};

			rminder.titleInput.value = '';
			rminder.detailsContainer.dataset.id = '2';

			rminder.renameTask(db);

			td.verify(db.postMessage(), { times: 0, ignoreExtraArgs: true });
			td.verify(console.log('Required field(s) missing: title', []), { times: 1 });
		});
	});

	describe('Rminder.showDetails', () => {
		it('Should send postMessage to show details', () => {
			td.replace(rminder, 'screenTest');
			td.replace(rminder, 'setSelected');

			const db = {
				postMessage: td.func()
			};

			const task = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			const data = {
				type: 'showDetails',
				id: 1
			};

			rminder.tasks(task);

			rminder.detailsContainer.dataset.id = '2';

			const listTask = document.querySelector('vp-list-task');

			listTask.task = '1';

			expect(rminder.detailsContainer.classList.contains('expanded')).to.be.false;

			rminder.showDetails(listTask, db);

			td.verify(db.postMessage(data));
			td.verify(rminder.screenTest());
			td.verify(rminder.setSelected(listTask));

			expect(rminder.detailsContainer.classList.contains('expanded')).to.be.true;
		});

		it('Should not send postMessage if details are already loaded', () => {
			td.replace(rminder, 'screenTest');
			const db = {
				postMessage: td.func()
			};

			const task = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(task);

			rminder.detailsContainer.dataset.id = '1';

			document.querySelector('vp-list-task').task = '1';

			expect(rminder.detailsContainer.classList.contains('expanded')).to.be.false;

			rminder.showDetails(document.querySelector('vp-list-task'), db);

			td.verify(db.postMessage(), { times: 0, ignoreExtraArgs: true });
			td.verify(rminder.screenTest());
			expect(rminder.detailsContainer.classList.contains('expanded')).to.be.true;
		});

		// it('Should not call postMessage and add expanded class is element does not have class "show-details"', () => {
		// 	td.replace(rminder, 'screenTest');
		// 	const db = {
		// 		postMessage: td.func()
		// 	};

		// 	expect(rminder.detailsContainer.classList.contains('expanded')).to.be.false;

		// 	rminder.showDetails(document.querySelector('.tasks'), db);

		// 	td.verify(db.postMessage(), { times: 0, ignoreExtraArgs: true });
		// 	td.verify(rminder.screenTest(), { times: 0 });
		// 	expect(rminder.detailsContainer.classList.contains('expanded')).to.be.false;
		// });
	});

	describe('Rminder.hideDetails', () => {
		it('Should hide details', () => {
			td.replace(rminder, 'screenTest');

			rminder.detailsContainer.classList.add('expanded');
			rminder.modal.classList.add('open');

			rminder.hideDetails();

			expect(rminder.detailsContainer.classList.contains('expanded')).to.be.false;
			expect(rminder.modal.classList.contains('open')).to.be.false;
			td.verify(rminder.screenTest());
		});
	});

	describe('Rminder.toggleSidebar', () => {
		it('Should toggle sidebar', () => {
			td.replace(rminder, 'screenTest');

			rminder.toggleSidebar();

			expect(rminder.sidebar.classList.contains('expanded')).to.be.false;

			rminder.toggleSidebar();

			expect(rminder.sidebar.classList.contains('expanded')).to.be.true;
			td.verify(rminder.screenTest(undefined, rminder.detailsContainer));
		});
	});

	describe('Rminder.hideSidebar', () => {
		it('Should hide sidebar', () => {
			td.replace(rminder, 'screenTest');

			expect(rminder.sidebar.classList.contains('expanded')).to.be.true;

			rminder.hideSidebar();

			expect(rminder.sidebar.classList.contains('expanded')).to.be.false;
			td.verify(rminder.screenTest());
		});
	});

	describe('Rminder.handleEvent', () => {
		it('Should call postMessage', () => {
			const db = {
				postMessage: td.func()
			};

			const data = {
				id: 1,
				type: 'importantTask',
				list: ''
			}

			rminder.handleEvent('importantTask', db, 1);

			td.verify(db.postMessage(data));
		});

		it('Should call postMessage with container id', () => {
			const db = {
				postMessage: td.func()
			};

			const data = {
				id: 1,
				type: 'importantTask',
				list: ''
			}

			rminder.detailsContainer.dataset.id = '1';

			rminder.handleEvent('importantTask', db);

			td.verify(db.postMessage(data));
		});
	});

	describe('Rminder.setTaskNote', () => {
		it('Should send postMessage to add note', () => {
			const db = {
				postMessage: td.func()
			};

			const data = {
				type: 'noteTask',
				id: 1,
				note: 'something',
				list: ''
			};

			rminder.note.value = 'something';
			rminder.detailsContainer.dataset.id = '1';

			rminder.setTaskNote(rminder.note, db);

			td.verify(db.postMessage(data));
		});
	});

	describe('Rminder.showList', () => {
		it('Should send postMessage to add note', () => {
			const db = {
				postMessage: td.func()
			};

			const data = {
				type: 'list',
				list: 'important'
			};

			rminder.showList('important', db);

			td.verify(db.postMessage(data));
		});
	});

	describe('Rminder.setDetailClasses', () => {
		it('Should set classes when data is an array', () => {
			const data = [{
				id: 1,
				completed: true,
				important: true,
				title: 'test',
				my_day: true,
				note: '',
				creation_date: 1616785736469
			}];

			rminder.detailsContainer.dataset.id = '1';

			rminder.setDetailClasses(data);

			expect(rminder.detailsContainer.classList.contains('important')).to.be.true;
			expect(rminder.detailsContainer.classList.contains('completed')).to.be.true;
			expect(rminder.detailsContainer.classList.contains('today')).to.be.true;
		});

		it('Should set classes when data is an object', () => {
			const data = {
				id: 1,
				completed: true,
				important: true,
				title: 'test',
				my_day: true,
				note: '',
				creation_date: 1616785736469
			};

			rminder.detailsContainer.dataset.id = '1';

			rminder.setDetailClasses(data);

			expect(rminder.detailsContainer.classList.contains('important')).to.be.true;
			expect(rminder.detailsContainer.classList.contains('completed')).to.be.true;
			expect(rminder.detailsContainer.classList.contains('today')).to.be.true;
		});
	});

	describe('Rminder.setDetailClasses', () => {
		it('Should remove hidden class from main container', () => {
			rminder.mainContainer.classList.add('hidden');

			rminder.screenTest();

			expect(rminder.mainContainer.classList.contains('hidden')).to.be.false;
		});

		it('Should add hidden class to main container if matches small media query', () => {
			rminder.smallMediaQuery = { matches: true };

			rminder.mainContainer.classList.remove('hidden');

			rminder.screenTest(rminder.mediaQueryList, rminder.sidebar);

			expect(rminder.mainContainer.classList.contains('hidden')).to.be.true;
		});

		it('Should remove sidebar expanded', () => {
			rminder.detailsContainer.classList.add('expanded');

			expect(rminder.sidebar.classList.contains('expanded')).to.be.true;

			rminder.screenTest(rminder.mediaQueryList, rminder.sidebar);

			expect(rminder.sidebar.classList.contains('expanded')).to.be.false;
		});
	});

	describe('Rminder.setSelected', () => {
		it('Should add selected class', () => {
			const task = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(task);

			const elem = document.querySelector('vp-list-task');

			expect(elem.classList.contains('selected')).to.be.false;

			rminder.setSelected(elem);

			expect(elem.classList.contains('selected')).to.be.true;
		});
	});

	describe('Rminder.settingsCompleted', () => {
		it('Should postMessage with completed true', () => {
			const db = {
				postMessage: td.func()
			};

			rminder.toggleCompleted.checked = true;

			rminder.settingsCompleted(rminder.toggleCompleted, db);

			td.verify(db.postMessage({ type: 'settings', completed: true, list: '' }));
		});
		it('Should postMessage with completed false', () => {
			const db = {
				postMessage: td.func()
			};

			rminder.toggleCompleted.checked = false;

			rminder.settingsCompleted(rminder.toggleCompleted, db);

			td.verify(db.postMessage({ type: 'settings', completed: false, list: '' }));
		});
	});

	describe('Rminder.toggle', () => {
		it('Should toggle open class', () => {
			rminder.settingsBtn.classList.remove('open');

			rminder.toggle(rminder.settingsBtn);

			expect(rminder.settingsBtn.classList.contains('open')).to.be.true;

			rminder.toggle(rminder.settingsBtn);

			expect(rminder.settingsBtn.classList.contains('open')).to.be.false;
		});
	});

	describe('Rminder.settings', () => {
		it('Should return false if no settings', () => {
			expect(rminder.settings()).to.be.false;
		});
		it('Should hide Completed tasks', () => {
			const data = {
				settings: {
					completed: 'hide'
				}
			}

			rminder.settings(data);

			expect(rminder.toggleCompleted.checked).to.be.true;
			expect(rminder.listCompleted.classList.contains('hidden')).to.be.true;
		});
		it('Should select Tasks list if Completed list is selected before hidding', () => {
			const data = {
				settings: {
					completed: 'hide'
				}
			}

			rminder.listCompleted.classList.add('selected');

			td.replace(rminder.listTasks, 'click')

			rminder.settings(data);

			expect(rminder.toggleCompleted.checked).to.be.true;
			expect(rminder.listCompleted.classList.contains('hidden')).to.be.true;
			td.verify(rminder.listTasks.click());
		});
		it('Should show Completed tasks', () => {
			const data = {
				settings: {
					completed: 'show'
				}
			};

			rminder.settings(data);

			expect(rminder.toggleCompleted.checked).to.be.false;
			expect(rminder.listCompleted.classList.contains('hidden')).to.be.false;
		});
		it('Should set order filter', () => {
			const data = {
				settings: {
					filter: 'important'
				}
			};

			expect(rminder.orderFilters[0].checked).to.be.false;

			rminder.settings(data);

			expect(rminder.orderFilters[0].checked).to.be.true;
		});
	});

	describe('Rminder.filterUpdate', () => {
		let db;

		beforeEach(() => {
			db = {
				postMessage: td.func()
			};
		});

		it('Should postMessage with filter important', () => {
			rminder.filterUpdate(rminder.orderFilters[0], db);

			td.verify(db.postMessage({ type: 'filter', filter: 'important' }));
		});

		it('Should postMessage with filter oldest', () => {
			rminder.filterUpdate(rminder.orderFilters[1], db);

			td.verify(db.postMessage({ type: 'filter', filter: 'oldest' }));
		});

		it('Should postMessage with filter newest', () => {
			rminder.filterUpdate(rminder.orderFilters[2], db);

			td.verify(db.postMessage({ type: 'filter', filter: 'newest' }));
		});
	});

	describe('Rminder.addEventListeners', () => {
		let db;

		beforeEach(() => {
			db = {
				postMessage: td.func()
			};
			rminder.addEventListeners(db);
		});

		afterEach(() => {
			td.reset();
		});

		it('Should call addTask', () => {
			td.replace(rminder, 'addTask');

			rminder.addTaskBtn.click();

			td.verify(rminder.addTask(db));
		});

		it('Should call handleEvent for myDayTask', () => {
			td.replace(rminder, 'handleEvent');

			rminder.myDay.click();

			td.verify(rminder.handleEvent('myDayTask', db));
		});

		it('Should call handleEvent for importantTask', () => {
			td.replace(rminder, 'handleEvent');

			rminder.importanceCheckBtn.click();

			td.verify(rminder.handleEvent('importantTask', db, 0));
		});

		it('Should call handleEvent for removeTask', () => {
			td.replace(rminder, 'handleEvent');

			const modalDelete = document.querySelector('vp-button.warning');
			modalDelete.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'removeTask'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.handleEvent('removeTask', db));
		});

		it('Should add open class to modal', () => {
			const remove = document.querySelector('.remove');
			remove.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'openModal'
					},
					bubbles: true,
					composed: true
				})
			);

			expect(rminder.modal.classList.contains('open')).to.be.true;
		});

		it('Should remove open class to modal', () => {
			rminder.modal.classList.add('open');

			const modalCancel = document.querySelector('vp-button.default');
			modalCancel.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'closeModal'
					},
					bubbles: true,
					composed: true
				})
			);

			expect(rminder.modal.classList.contains('open')).to.be.false;
		});

		it('Should call handleEvent for completedTask', () => {
			td.replace(rminder, 'handleEvent');
			td.replace(rminder, 'hideDetails');

			rminder.completedCheck.click();

			td.verify(rminder.handleEvent('completedTask', db, 0));
			td.verify(rminder.hideDetails(), { times: 0 });
		});

		it('Should call handleEvent for completedTask and hideDetails', () => {
			td.replace(rminder, 'handleEvent');
			td.replace(rminder, 'hideDetails');
			rminder.toggleCompleted.checked = true;

			rminder.completedCheck.click();

			td.verify(rminder.handleEvent('completedTask', db, 0));
			td.verify(rminder.hideDetails());
		});

		it('Should call setTaskNote', () => {
			td.replace(rminder, 'setTaskNote');

			const noteBtn = document.querySelector('.add-note');
			noteBtn.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'setTaskNote'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.setTaskNote(noteBtn, db));
		});

		it('Should call renameTask', () => {
			td.replace(rminder, 'renameTask');

			const rename = document.querySelector('.rename');
			rename.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'renameTask'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.renameTask(rename, db));
		});

		it('Should call renameTask on keyup', () => {
			td.replace(rminder, 'renameTask');

			rminder.titleInput.dispatchEvent(new KeyboardEvent('keyup', { key: 'Enter', code: 'Enter' }));

			td.verify(rminder.renameTask(db));
		});

		it('Should not call renameTask on keyup if key is not enter', () => {
			td.replace(rminder, 'renameTask');

			rminder.titleInput.dispatchEvent(new KeyboardEvent('keyup', { key: 'a', code: 'KeyA' }));

			td.verify(rminder.renameTask(), { times: 0, ignoreExtraArgs: true });
		});

		it('Should call addTask on keyup', () => {
			td.replace(rminder, 'addTask');

			rminder.taskInput.dispatchEvent(new KeyboardEvent('keyup', { key: 'Enter', code: 'Enter' }));

			td.verify(rminder.addTask(db));
		});

		it('Should not call addTask on keyup if key is not enter', () => {
			td.replace(rminder, 'addTask');

			rminder.taskInput.dispatchEvent(new KeyboardEvent('keyup', { key: 'a', code: 'KeyA' }));

			td.verify(rminder.addTask(), { times: 0, ignoreExtraArgs: true });
		});

		it('Should call showList on list click', () => {
			rminder.smallMediaQuery = { matches: true };
			td.replace(rminder, 'showList');
			td.replace(rminder, 'hideSidebar');

			rminder.listMyDay.dispatchEvent(
				new CustomEvent('vp-nav-list:click', {
					detail: {
						name: 'my_day'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.showList('my_day', db));
			td.verify(rminder.hideSidebar());
		});

		it('Should not call hideSidebar on list click', () => {
			td.replace(rminder, 'showList');
			td.replace(rminder, 'hideSidebar');

			rminder.listMyDay.dispatchEvent(
				new CustomEvent('vp-nav-list:click', {
					detail: {
						name: 'my_day'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.showList('my_day', db));
			td.verify(rminder.hideSidebar(), { times: 0 });
		});

		it('Should not call showList and hideSidebar on list click if list does not have name', () => {
			td.replace(rminder, 'showList');
			td.replace(rminder, 'hideSidebar');

			rminder.listMyDay.dispatchEvent(
				new CustomEvent('vp-nav-list:click', {
					detail: {},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.showList(), { times: 0, ignoreExtraArgs: true });
			td.verify(rminder.hideSidebar(), { times: 0 });
		});

		it('Should not call showList', () => {
			td.replace(rminder, 'showList');

			rminder.lists.click();

			td.verify(rminder.showList(), { times: 0, ignoreExtraArgs: true });
		});

		it('Should call showDetails on task click', () => {
			rminder.smallMediaQuery = { matches: true };
			td.replace(rminder, 'showDetails');

			const task = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(task);

			const listTask = document.querySelector('vp-list-task');
			listTask.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'showDetails'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.showDetails(listTask, db));
		});

		it('Should not call showDetails on task click', () => {
			rminder.smallMediaQuery = { matches: true };
			td.replace(rminder, 'showDetails');

			rminder.taskList.click();

			td.verify(rminder.showDetails(), { times: 0, ignoreExtraArgs: true });
		});

		it('Should call handleEvent with importantTask', () => {
			td.replace(rminder, 'handleEvent');
			const task = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(task);

			const listTask = document.querySelector('vp-list-task');
			listTask.task = '1';

			listTask.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'importantTask'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.handleEvent('importantTask', db, 1));
		});

		it('Should call handleEvent with completedTask', () => {
			td.replace(rminder, 'handleEvent');
			td.replace(rminder, 'hideDetails');
			const task = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(task);

			const listTask = document.querySelector('vp-list-task');
			listTask.task = '1';

			listTask.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'completedTask'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.handleEvent('completedTask', db, 1));
			td.verify(rminder.hideDetails(), { times: 0 });
		});

		it('Should call handleEvent with completedTask and hideDetails', () => {
			td.replace(rminder, 'handleEvent');
			td.replace(rminder, 'hideDetails');
			rminder.toggleCompleted.checked = true;
			rminder.detailsContainer.dataset.id = '1';
			const task = {
				value: [{
					id: 1,
					completed: true,
					important: true,
					title: 'test',
					my_day: false,
					note: '',
					creation_date: 1616785736469
				}]
			};

			rminder.tasks(task);

			const listTask = document.querySelector('vp-list-task');
			listTask.task = '1';

			listTask.dispatchEvent(
				new CustomEvent('vp-button:click', {
					detail: {
						trigger: 'completedTask'
					},
					bubbles: true,
					composed: true
				})
			);

			td.verify(rminder.handleEvent('completedTask', db, 1));
			td.verify(rminder.hideDetails());
		});

		it('Should call toggle of filter button', () => {
			td.replace(rminder, 'toggle');

			rminder.filterBtn.click();

			td.verify(rminder.toggle(rminder.filterBtn));
		});

		it('Should call toggle of settings button', () => {
			td.replace(rminder, 'toggle');

			rminder.settingsBtn.click();

			td.verify(rminder.toggle(rminder.settingsBtn));
		});

		it('Should call settingsCompleted of toggleCompleted button', () => {
			td.replace(rminder, 'settingsCompleted');

			rminder.toggleCompleted.click();

			td.verify(rminder.settingsCompleted(rminder.toggleCompleted, db));
		});

		it('Should call filterUpdate on change of order filter', () => {
			td.replace(rminder, 'filterUpdate');

			rminder.orderFilters[0].click();

			td.verify(rminder.filterUpdate(rminder.orderFilters[0], db));
		});
	});
});
