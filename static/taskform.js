const RECURRING_INPUT_UNIT_SELECTOR = "#task-form-recurring-unit";
const TASK_FORM_SELECTOR = "#task-form"

const updateRecurringUnits = (event) => {
  if (event == null || event.target == null) return;

  let recurringValue = event.target.value;
  let recurringUnit = document.querySelector(RECURRING_INPUT_UNIT_SELECTOR);
  if (recurringUnit == null) {
    console.log(
      "unable to locate recurring unit input using " +
        RECURRING_INPUT_UNIT_SELECTOR,
    );
    return;
  }
  let options = recurringUnit.children;
  if (options.length < 1) {
    console.log(
      "unable to find <option> elements within the recurring unit input",
    );
    return;
  }

  if (recurringValue == 1) {
    for (const opt of options) {
      if (opt.textContent.at(-1) == "s")
        opt.textContent = opt.textContent.slice(0, -1);
    }
  } else {
    for (const opt of options) {
      if (opt.textContent.at(-1) != "s")
        opt.textContent = opt.textContent + "s";
    }
  }
};

const enterSubmitForm = (event) => {
  if (event.key === "Enter" && !event.shiftKey) {
    let form = document.querySelector(TASK_FORM_SELECTOR);
    if (form == null) {
      console.log("couldn't find a task form using " + TASK_FORM_SELECTOR);
      return;
    }

    event.preventDefault();
    form.requestSubmit();
  }
};
