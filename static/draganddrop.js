const draggableClassName = "dashboard-list-item-container";
const containerClassName = "dashboard-droppable-container";
const hoveredDroppableContainerClass = "dashboard-hovered-droppable-container";
const categoriesContainerClassName = "dashboard-categories-container";

let draggedItem = null;

findParentContainer = function (el) {
  let parent = el.parentElement;
  while (parent != null && !parent.classList.contains(containerClassName)) {
    parent = parent.parentElement;
  }
  return parent;
};

class DraggableItem {
  constructor(element, event, id) {
    let rect = element.getBoundingClientRect();

    this.el = element;
    this.id = id;
    this.category =
      element.parentElement.parentElement.parentElement.dataset.category;
    this.offsetX = event.clientX - rect.left;
    this.offsetY = event.clientY - rect.top;
    this.parent = element.parentElement;
    this.nextSib = element.nextElementSibling;
    this.parentContainer = findParentContainer(element);
    this.moving = false;
  }

  updatePosition(x, y) {
    this.x = x;
    this.y = y;
    this.el.style.top = this.y - this.offsetY + "px";
    this.el.style.left = this.x - this.offsetX + "px";
  }

  initStyles() {
    this.el.style.pointerEvents = "none";
    this.el.style.position = "absolute";
    document.body.style.userSelect = "none";
    // document.body.style.position = "relative";
  }

  revertStyles() {
    this.el.style.pointerEvents = "auto";
    this.el.style.position = "";
    document.body.style.userSelect = "";
    // document.body.style.position = "";
  }
}

resetHoveredDraggableContainer = function () {
  let droppableContainers = document.querySelectorAll("." + containerClassName);
  for (const cont of droppableContainers) {
    cont.classList.remove(hoveredDroppableContainerClass);
  }
};

onMouseDown = (e, el) => {
  let id = el.firstChild.dataset.id;
  if (id === "") return;

  draggedItem = new DraggableItem(el, e, id);
};

document.body.addEventListener("mousemove", (e) => {
  if (draggedItem == null || draggedItem.el == null) return;

  if (!draggedItem.moving) {
    draggedItem.initStyles();
    document.body.appendChild(draggedItem.el);
    draggedItem.moving = true;
  }
  draggedItem.updatePosition(e.clientX, e.clientY);
  let hoveredContainer = document
    .elementFromPoint(e.clientX, e.clientY)
    .closest("." + containerClassName);
  if (hoveredContainer == null) {
    resetHoveredDraggableContainer();
    draggedItem.date = "";
  } else {
    resetHoveredDraggableContainer();
    hoveredContainer.classList.add(hoveredDroppableContainerClass);
    draggedItem.date = hoveredContainer.id.replace("date-", "");
  }
});

document.body.addEventListener("mouseup", (_) => {
  if (draggedItem == null || draggedItem.el == null) return;

  if (draggedItem.x == null || draggedItem.y == null) {
    draggedItem = null;
    return;
  }

  let hoveredContainer = document
    .elementFromPoint(draggedItem.x, draggedItem.y)
    .closest("." + containerClassName);

  if (
    hoveredContainer == null ||
    hoveredContainer === draggedItem.parentContainer
  ) {
    resetHoveredDraggableContainer();
    draggedItem.parent.insertBefore(draggedItem.el, draggedItem.nextSib);
  } else {
    let endpoint = "/update-task-duedate?id=" + draggedItem.id;
    if (draggedItem.date !== "dashboard-unscheduled-tasks") {
      endpoint = endpoint + "&date=" + draggedItem.date;
    }
    htmx.ajax("GET", endpoint, {
      target: "#" + hoveredContainer.id,
      swap: "outerHTML",
    });
    let categoriesContainer = draggedItem.parent;
    if (categoriesContainer.children.length == 0) {
      while (
        !categoriesContainer.classList.contains(categoriesContainerClassName)
      ) {
        categoriesContainer = categoriesContainer.parentElement;
      }
      for (const c of categoriesContainer.children) c.remove();
    }
    draggedItem.el.remove();
  }

  draggedItem.revertStyles();
  draggedItem = null;
});

// document
// .querySelectorAll("." + draggableClassName)
// .forEach((el) => makeDraggable(el));
