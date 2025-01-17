setTimeout(postCode, 200);

async function postCode() {
  const outputArea = document.querySelector("#output-area");
  const errorArea = document.querySelector("#error-area");
  const errorContainer = document.querySelector("#error-container");
  const argument = document.querySelector("#argument-input").value;
  const functionDefinition = document.querySelector(
    "#function-def-input"
  ).value;

  try {
    const res = await fetch("/run", {
      method: "POST",
      body: JSON.stringify({ functionDefinition, argument }),
    });
    if (!res.ok) {
      throw new Error(await res.text());
    }
    const output = await res.text();

    if (
      argument !== document.querySelector("#argument-input").value ||
      functionDefinition !== document.querySelector("#function-def-input").value
    ) {
      // The inputs have changed since this request was sent; ignore the result.
      return;
    }
    errorContainer.style.display = "none";
    errorArea.innerHTML = "";
    outputArea.innerHTML = output;
  } catch (e) {
    errorContainer.style.display = "block";
    errorArea.innerHTML = e.message;
    outputArea.innerHTML = "";
  }
}
