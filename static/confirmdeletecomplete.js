const DELETE_TASK_HX_TARGET = ".task-info-content";
const COMPLETE_TASK_HX_TARGET = ".task-info-content";
const CONFIRM_DELETE_BUTTON_SELECTOR = "#task-info-confirm-delete-button";
const CONFIRM_COMPLETE_BUTTON_SELECTOR = "#task-info-confirm-complete-button";

const confirmDelete = (event) => {
  let deleteButton = document.createElement("div");
  deleteButton.classList.add("task-info-delete-button", "div-button");

  let id = event.target.dataset.id;
  if (id == null) {
    console.log("no id set on");
    console.log(event.target);
    return;
  }

  deleteButton.setAttribute("hx-get", "/delete?id=" + id);
  deleteButton.setAttribute("hx-target", DELETE_TASK_HX_TARGET);
  deleteButton.setAttribute("hx-swap", "outerHTML");
  deleteButton.innerHTML = "Confirm Delete?";

  let confirmDeleteButton = document.querySelector(CONFIRM_DELETE_BUTTON_SELECTOR);
  if (confirmDeleteButton == null) {
    console.log(
      "cannot locate confirm delete button using" +
        CONFIRM_DELETE_BUTTON_SELECTOR,
    );
    return;
  }

  let nextSib = confirmDeleteButton.nextSibling;
  let parent = confirmDeleteButton.parentElement;
  if (parent == null) {
    console.log("error locating parentElement for");
    console.log(confirmDeleteButton);
    return;
  }

  parent.insertBefore(deleteButton, nextSib);
  htmx.process(deleteButton);

  confirmDeleteButton.remove();
};

const confirmComplete = (event) => {
  let completeButton = document.createElement("div");
  completeButton.classList.add("task-info-complete-button", "div-button");

  let id = event.target.dataset.id
  if (id == null) {
    console.log("no id set on");
    console.log(event.target);
    return;
  }

  completeButton.setAttribute("hx-get", "/complete?id=" + id);
  completeButton.setAttribute("hx-target", COMPLETE_TASK_HX_TARGET);
  completeButton.setAttribute("hx-swap", "outerHTML");
  completeButton.innerHTML = "Confirm Complete?";

  let confirmCompleteButton = document.querySelector(CONFIRM_COMPLETE_BUTTON_SELECTOR);
  if (confirmCompleteButton == null) {
    console.log(
      "cannot locate confirm complete button using" +
        CONFIRM_COMPLETE_BUTTON_SELECTOR,
    );
    return;
  }

  let nextSib = confirmCompleteButton.nextSibling;
  let parent = confirmCompleteButton.parentElement;
  if (parent == null) {
    console.log("error locating parentElement for");
    console.log(confirmCompleteButton);
    return;
  }

  parent.insertBefore(completeButton, nextSib);
  htmx.process(completeButton);

  confirmCompleteButton.remove();
}
