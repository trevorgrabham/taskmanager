const removeCompleteTaskTimeout = 5000;

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


