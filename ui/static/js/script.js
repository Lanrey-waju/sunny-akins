document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("contactForm");
  const button = document.getElementById("showContactForm");

  // Toggle form visibility
  button?.addEventListener("click", () => {
    form.classList.toggle("show");
    button.textContent = form.classList.contains("show")
      ? "Hide Contact Form"
      : "Request a Quote";
  });

  // Clear error state for an input
  const clearError = (input) => {
    const errorSpan = input.parentElement.querySelector(".error-message");
    if (errorSpan) {
      errorSpan.remove();
    }
    input.classList.remove("border-red-500");
  };

  // Show error for an input
  const showError = (input, message) => {
    clearError(input);
    input.classList.add("border-red-500");
    const errorSpan = document.createElement("span");
    errorSpan.className = "error-message text-red-500 text-sm mt-1";
    errorSpan.textContent = message;
    input.parentElement.appendChild(errorSpan);
  };

  // Clear all errors
  const clearAllErrors = () => {
    form.querySelectorAll("input, textarea").forEach(clearError);
  };

  // Form submission handler
  window.submitForm = function (event) {
    event.preventDefault();
    clearAllErrors();

    const formData = new FormData(form);
    const submitButton = form.querySelector('button[type="submit"]');

    // Disable submit button
    submitButton.disabled = true;
    submitButton.classList.add("opacity-50");

    fetch("/contact/post", {
      method: "POST",
      body: formData,
      headers: {
        Accept: "application/json",
      },
    })
      .then(async (response) => {
        const data = await response.json();

        if (response.ok) {
          // Success case
          alert("Message sent successfully!");
          form.reset();
          form.classList.remove("show");
          button.textContent = "Request a Quote";
          window.location.href = "/";
        } else if (response.status === 422) {
          // Validation errors
          if (data.errors) {
            // Display each error
            Object.entries(data.errors).forEach(([field, message]) => {
              const input = form.querySelector(`[name="${field}"]`);
              if (input) {
                showError(input, message);
              }
            });

            // Repopulate form if needed
            if (data.formData) {
              Object.entries(data.formData).forEach(([field, value]) => {
                const input = form.querySelector(`[name="${field}"]`);
                if (input) {
                  input.value = value;
                }
              });
            }
          }
        } else {
          throw new Error("Form submission failed");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("There was an error sending your message. Please try again.");
      })
      .finally(() => {
        // Re-enable submit button
        submitButton.disabled = false;
        submitButton.classList.remove("opacity-50");
      });
  };

  // Clear errors when user starts typing
  form.querySelectorAll("input, textarea").forEach((input) => {
    input.addEventListener("input", () => clearError(input));
  });
});
