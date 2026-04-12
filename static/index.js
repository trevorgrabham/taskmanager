const removeCompleteTaskTimeout = 5000;
const selectedSuggestionClassName = "active";
const categoryMatchClassName = "add-form-category-matched";
const maxCategorySelections = 5;

let categorySuggestionIndex = -1;
let categoryInputTimeout;

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
