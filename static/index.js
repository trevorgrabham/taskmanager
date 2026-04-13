const removeCompleteTaskTimeout = 5000;
const selectedSuggestionClassName = "active";
const categoryMatchClassName = "add-form-category-matched";
const maxCategorySelections = 5;

let categorySuggestionIndex = -1;
let categoryInputTimeout;
let draggedTask = null;

function removeCompletedTask(taskEl) {
  let parent;
  setTimeout(() => {
    if (taskEl.classList.contains("completed-task")) {
      parent = taskEl.parentNode.parentNode;
      taskEl.parentNode.remove();
      setTimeout(() => {
        if (parent.children.length <= 1) {
          parent.remove();
        }
      }, 0);
    }
  }, removeCompleteTaskTimeout);
}

function getCategorySuggestions() {
  return [...document.querySelectorAll(".add-form-category-suggestion")];
}

function getMatchedCategorySuggestions() {
  return [
    ...document.querySelectorAll(
      ".add-form-category-suggestion." + categoryMatchClassName,
    ),
  ];
}

function clearCategorySelection() {
  console.log("clearing category selection");

  getCategorySuggestions().forEach((el) =>
    el.classList.remove(selectedSuggestionClassName),
  );
}

function clearCategoryMatches() {
  console.log("clearing category matches");

  getCategorySuggestions().forEach((el) =>
    el.classList.remove(categoryMatchClassName),
  );
}

function selectCategorySuggestion(index) {
  console.log("selecting category suggestion for", index);

  const suggestions = getMatchedCategorySuggestions();
  if (suggestions.length <= 0) {
    return;
  }

  clearCategorySelection();
  categorySuggestionIndex = index;
  if (index < 0) {
    categorySuggestionIndex = suggestions.length - 1;
  }
  if (index >= suggestions.length) {
    categorySuggestionIndex = 0;
  }
  suggestions[categorySuggestionIndex].classList.add(
    selectedSuggestionClassName,
  );
}

function applyCategorySelection() {
  console.log("applying category suggestion");

  const suggestions = getMatchedCategorySuggestions();
  const categoryInput = document.getElementById("add-form-category");
  if (categoryInput == null) {
    return;
  }
  if (
    categorySuggestionIndex >= 0 &&
    categorySuggestionIndex < suggestions.length
  ) {
    categoryInput.value = suggestions[categorySuggestionIndex].innerHTML;
  }
}

function filterCategorySuggestions(inputValue) {
  console.log("filtering category suggestion");

  const suggestions = getCategorySuggestions();
  if (suggestions.length <= 0) {
    return;
  }

  clearCategorySelection();

  suggestions.forEach((suggestion) => {
    suggestion.classList.remove("add-form-category-matched");
  });

  suggestions
    .filter((suggestion) => {
      return suggestion.innerHTML.includes(inputValue);
    })
    .slice(0, maxCategorySelections)
    .forEach((suggestion) => {
      suggestion.classList.add("add-form-category-matched");
    });
}

function registerCategorySuggestions() {
  console.log("registering suggestion listeners");

  const categoryInput = document.getElementById("add-form-category");
  if (categoryInput == null) {
    return;
  }
  const form = document.getElementById("add-form");
  if (form == null) {
    return;
  }

  categoryInput.addEventListener("input", (e) => {
    clearTimeout(categoryInputTimeout);
    categorySuggestionIndex = -1;

    categoryInputTimeout = setTimeout(() => {
      filterCategorySuggestions(e.target.value);
    }, 100);
  });

  categoryInput.addEventListener("keydown", (e) => {
    const suggestions = getMatchedCategorySuggestions();
    if (suggestions.length <= 0) {
      return;
    }

    if (e.key === "ArrowDown") {
      e.preventDefault();
      selectCategorySuggestion(categorySuggestionIndex + 1);
    }

    if (e.key === "ArrowUp") {
      e.preventDefault();
      selectCategorySuggestion(categorySuggestionIndex - 1);
    }

    if (e.key === "Tab" && categorySuggestionIndex !== -2) {
      if (e.shiftKey) {
        e.preventDefault();
        selectCategorySuggestion(categorySuggestionIndex - 1);
      } else {
        e.preventDefault();
        selectCategorySuggestion(categorySuggestionIndex + 1);
      }
    }

    if (e.key === "Enter") {
      e.preventDefault();

      if (
        categorySuggestionIndex >= 0 &&
        categorySuggestionIndex < suggestions.length
      ) {
        applyCategorySelection();
      }

      clearCategorySelection();
      clearCategoryMatches();
      categorySuggestionIndex = -2;
      let nextSib = categoryInput.nextSibling;
      while (nextSib.tagName !== "INPUT") {
        nextSib = nextSib.nextSibling;
      }
      nextSib.focus();
    }
  });

  form.addEventListener("click", (e) => {
    const categoryInput = document.getElementById("add-form-category");
    if (categoryInput == null) {
      return;
    }

    if (e.target.classList.contains("add-form-category-suggestion")) {
      categoryInput.value = e.target.innerHTML;
    }
  });
}

// func registerDashboardForDrags() {
//   const dashboard = document.getElementById("dashboard")
//   if (dashboard == null) { return }
//
//   dashboard.addEventListener("mousemove", (e) => {
//     if (draggedTask == null || draggedTask.task == null) { return }
//
//     // update the x and y values for draggedTask 
//     //    clientX/Y - getBoundingClientRect().left/top
//     // change its positioning. style.top, style.left
//   })
//
//   dashboard.addEventListener("mouseup", (e) => {
//     if (draggedTask == null || draggedTask.task == null) { return }
//
//     // get the element underneath the current mouse position. document.elementFromPoint(x, y).closest("droppable_container")
//     // if there is no droppable container 
//     //    revert back to previous parent. parent.insertBefore(draggedTask, nextSib)
//     // else
//     //    see if the droppable container has a matching category header
//     //    if they do 
//     //      iterate through its children until we find a child that has a title that comes after draggedTask
//     //      if we reach the end of the children then just parent.appendChild(draggedTask)
//     //    else 
//     //      hit an endpoint that sends the taskID and containers date so create a new category
//     // reset pointerevents. style.pointerEvents = "auto"
//     // draggedTask = null
//   })
// }
//
// func registerDraggable(el) {
//   if (el == null) { return }
//
//   el.addEventListener("mousedown", (e) => {
//     // grab the parent and nextSibling from draggedTask 
//     // grab the category
//     // move draggedTask to a child of the <body> instead of its parent doc.body.appendChild(draggedEl)
//     // remove pointer events from the draggedTask .style.pointerEvents = "none"
//   })
// }
