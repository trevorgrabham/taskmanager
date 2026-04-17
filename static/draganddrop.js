const draggableClassName = "dashboard-list-item";
const taskClassName = "dashboard-list-item-container";
const taskParentClassName = "dashboard-list-items-container";
const droppableContainerClassName = "dashboard-droppable-container";
const hoveredDroppableContainerClassName =
  "dashboard-hovered-droppable-container";
const categoryContainerClassName = "dashboard-category-container";

let draggedTask;

class DraggableItem {
  constructor(element, event, id) {
    this.el = element.closest("." + taskClassName);
    this.moving = false;

    // for reverting back to previous spot
    this.parent = element.closest("." + taskParentClassName);
    this.nextSib = this.el.nextSibling;

    // for styling
    let rect = element.getBoundingClientRect();
    this.offsetX = event.clientX - rect.left;
    this.offsetY = event.clientY - rect.top;
    this.originalDroppableContainer = element.closest(
      "." + droppableContainerClassName,
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
}

// on mouse down
//
const registerDraggable = (event) => {
  if (!event.target.classList.contains(draggableClassName)) return;

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
    console.log("started moving");

    draggedTask.moving = true;
    document.body.appendChild(draggedTask.el);

    console.log("added dragged task to body");

    if (draggedTask.parent.children.length <= 0) {
      console.log("parent container is now empty");

      draggedTask.categoryContainer = draggedTask.parent.closest(
        "." + categoryContainerClassName,
      );
      draggedTask.categoryContainerParent =
        draggedTask.categoryContainer.parentElement;
      draggedTask.categoryContainerSib =
        draggedTask.categoryContainer.nextSibling;
      draggedTask.categoryContainer.remove();
    }
    draggedTask.initStyles();
    if (draggedTask.hoveredDroppableContainer != null) {
      draggedTask.hoveredDroppableContainer.classList.add(
        hoveredDroppableContainerClassName,
      );
    }
  }

  draggedTask.updatePosition(event);
  let hoveredDroppableContainer = document
    .elementFromPoint(event.clientX, event.clientY)
    ?.closest("." + droppableContainerClassName);

  if (hoveredDroppableContainer == null) {
    if (draggedTask.hoveredDroppableContainer != null)
      draggedTask.hoveredDroppableContainer.classList.remove(
        hoveredDroppableContainerClassName,
      );
    draggedTask.hoveredDroppableContainer = null;
  } else if (
    hoveredDroppableContainer !== draggedTask.hoveredDroppableContainer
  ) {
    if (draggedTask.hoveredDroppableContainer != null)
      draggedTask.hoveredDroppableContainer.classList.remove(
        hoveredDroppableContainerClassName,
      );
    hoveredDroppableContainer.classList.add(hoveredDroppableContainerClassName);
    draggedTask.hoveredDroppableContainer = hoveredDroppableContainer;
  }
});

document.addEventListener("mouseup", (event) => {
  if (draggedTask == null) return;
  if (!draggedTask.moving) {
    draggedTask = null;
    return;
  }

  console.log("mouse up event triggered");

  draggedTask.hoveredDroppableContainer?.classList.remove(
    hoveredDroppableContainerClassName,
  );
  hoveredDroppableContainer = document
    .elementFromPoint(event.clientX, event.clientY)
    .closest("." + droppableContainerClassName);
  if (
    hoveredDroppableContainer == null ||
    hoveredDroppableContainer === draggedTask.originalDroppableContainer
  ) {
    draggedTask.revertStyles();

    console.log("dropped outside of a container or back in original container");

    if (draggedTask.categoryContainer != null) {
      console.log("re-adding original category container");

      draggedTask.categoryContainerParent.insertBefore(
        draggedTask.categoryContainer,
        draggedTask.categoryContainerSib,
      );
    }

    draggedTask.parent.insertBefore(draggedTask.el, draggedTask.nextSib);
    draggedTask = null;
    return;
  }
  let endpoint = "/update-task-duedate?id=" + draggedTask.id;
  duedate = hoveredDroppableContainer.id;
  if (duedate == null) {
    console.log(
      "dropping dragged item: no id on ." +
        droppableContainerClassName +
        " container",
    );
    draggedTask.revertStyles();

    if (draggedTask.categoryContainer != null) {
      console.log("re-adding original category container");

      draggedTask.categoryContainerParent.insertBefore(
        draggedTask.categoryContainer,
        draggedTask.categoryContainerSib,
      );
    }
    draggedTask.parent.insertBefore(draggedTask.el, draggedTask.nextSib);
    draggedTask = null;
    return;
  } else if (duedate === "dashboard-unscheduled-tasks") {
    console.log("dropped into unscheduled");

    // do nothing
  } else {
    duedate = duedate.replace("date-", "");
    endpoint = endpoint + "&duedate=" + duedate;

    console.log("dropped into " + duedate);
  }

  htmx.ajax("GET", endpoint, {
    target: "#" + hoveredDroppableContainer.id,
    swap: "outerHTML",
  });
  draggedTask.el.remove();
  draggedTask = null;
});

document.addEventListener("mousedown", registerDraggable);
