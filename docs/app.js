import{Rminder as e}from"./js/Rminder.js";import{dbMessage as s}from"./js/dbMessage.js";(()=>{const r=new e;const a=new Worker("./js/workers/dbw.js");a.postMessage({type:"launch"});a.onmessage=e=>s(r,e.data,a);if(r.smallMediaQuery.matches){r.sidebar.classList.remove("expanded")}r.screenTest();r.setDocHeight()})();