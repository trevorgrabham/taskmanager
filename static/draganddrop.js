const DRAGGABLE_CLASS_NAME = "dashboard-list-item";
const TASK_CLASS_NAME = "dashboard-list-item-container";
const TASK_PARENT_CLASS_NAME = "dashboard-list-items-container";
const DROPPABLE_CONTAINER_CLASS_NAME = "dashboard-droppable-container";
const HOVERED_DROPPABLE_CONTAINER_CLASS_NAME =
  "dashboard-hovered-droppable-container";
const CATEGORY_CONTAINER_CLASS_NAME = "dashboard-category-container";

let draggedTask;

class DraggableItem {
  constructor(element, event, id) {
    this.el = element.closest("." + TASK_CLASS_NAME);
    this.moving = false;

    // for reverting back to previous spot
    this.parent = element.closest("." + TASK_PARENT_CLASS_NAME);
    this.nextSib = this.el.nextSibling;

    // for styling
    let rect = element.getBoundingClientRect();
    this.offsetX = event.clientX - rect.left;
    this.offsetY = event.clientY - rect.top;
    this.originalDroppableContainer = element.closest(
      "." + DROPPABLE_CONTAINER_CLASS_NAME,
    );
    this.hoveredDroppableContainer = this.originalDroppableContainer;
    // for backend database
    this.id = id;
  }

  updatePosition(event) {
    this.el.style.left = event.pageX - this.offsetX + "px";
    this.el.style.top = event.pageY - this.offsetY + "px";
  }

  initStyles() {
    this.el.style.pointerEvents = "none";
    this.el.style.position = "absolute";
    document.body.style.userSelect = "none";
  }

  revertStyles() {
    this.el.style.pointerEvents = "auto";
    this.el.style.position = "";
    document.body.style.userSelect = "";
  }

  trackCategoryContainer() {
    this.categoryContainer = this.parent.closest(
      "." + CATEGORY_CONTAINER_CLASS_NAME,
    );
    this.categoryContainerParent = this.categoryContainer.parentElement;
    this.categoryContainerSib = this.categoryContainer.nextSibling;
  }

  removeHoveredDroppableContainer() {
    if (this.hoveredDroppableContainer == null) return;

    this.hoveredDroppableContainer.classList.remove(
      HOVERED_DROPPABLE_CONTAINER_CLASS_NAME,
    );
  }

  revertCategoryAndTask() {
    if (this.categoryContainer != null) {
      console.log("re-adding original category container");

      this.categoryContainerParent.insertBefore(
        this.categoryContainer,
        this.categoryContainerSib,
      );
    }

    this.parent.insertBefore(this.el, this.nextSib);
  }
}

// on mouse down
//
const registerDraggable = (event) => {
  if (!event.target.classList.contains(DRAGGABLE_CLASS_NAME)) return;

  console.log("mouse down event triggered. target is");
  console.log(event.target);

  let taskID = event.target.dataset.id;
  if (taskID == null) {
    console.log("registering draggable: no 'data-id' attribute");
    return;
  }

  draggedTask = new DraggableItem(event.target, event, taskID);
};

document.addEventListener("mousemove", (event) => {
  if (draggedTask == null) return;

  console.log("mouse move event triggered");

  if (!draggedTask.moving) {
    draggedTask.moving = true;
    document.body.appendChild(draggedTask.el);

    if (draggedTask.parent.children.length <= 0) {
      draggedTask.trackCategoryContainer();
      draggedTask.categoryContainer.remove();
    }
    draggedTask.initStyles();
    if (draggedTask.hoveredDroppableContainer != null) {
      draggedTask.hoveredDroppableContainer.classList.add(
        HOVERED_DROPPABLE_CONTAINER_CLASS_NAME,
      );
    }
  }

  draggedTask.updatePosition(event);
  let hoveredDroppableContainer = document
    .elementFromPoint(event.clientX, event.clientY)
    ?.closest("." + DROPPABLE_CONTAINER_CLASS_NAME);

  if (hoveredDroppableContainer == null) {
    draggedTask.removeHoveredDroppableContainer();
    draggedTask.hoveredDroppableContainer = null;
  } else if (
    hoveredDroppableContainer !== draggedTask.hoveredDroppableContainer
  ) {
    draggedTask.removeHoveredDroppableContainer();
    hoveredDroppableContainer.classList.add(
      HOVERED_DROPPABLE_CONTAINER_CLASS_NAME,
    );
    draggedTask.hoveredDroppableContainer = hoveredDroppableContainer;
  }
});

document.addEventListener("mouseup", (event) => {
  if (draggedTask == null) return;
  if (!draggedTask.moving) {
    draggedTask = null;
    return;
  }

  draggedTask.removeHoveredDroppableContainer();
  hoveredDroppableContainer = document
    .elementFromPoint(event.clientX, event.clientY)
    ?.closest("." + DROPPABLE_CONTAINER_CLASS_NAME);
  if (
    hoveredDroppableContainer == null ||
    hoveredDroppableContainer === draggedTask.originalDroppableContainer
  ) {
    draggedTask.revertStyles();
    draggedTask.revertCategoryAndTask();
    draggedTask = null;
    return;
  }

  let endpoint = "/update-task-duedate?id=" + draggedTask.id;
  duedate = hoveredDroppableContainer.id;
  if (duedate == null) {
    console.log(
      "dropping dragged item: no data-id on " + draggedTask.el + " container",
    );
    draggedTask.revertStyles();
    draggedTask.revertCategoryAndTask();
    draggedTask = null;
    return;
  }

  if (/^date-\d{4}-\d{2}-\d{2}$/.test(duedate)) {
    endpoint = endpoint + "&duedate=" + duedate.replace("date-", "");
  }

  htmx.ajax("GET", endpoint, {
    target: "#" + hoveredDroppableContainer.id,
    swap: "outerHTML",
  });
  draggedTask.el.remove();
  draggedTask = null;
});

document.addEventListener("mousedown", registerDraggable);
