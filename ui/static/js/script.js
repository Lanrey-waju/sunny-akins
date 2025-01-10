// JavaScript for handling form visibility and submission
document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("contactForm");
  const button = document.getElementById("showContactForm");

  // Single event listener for toggle with proper transition handling
  button?.addEventListener("click", () => {
    if (form.classList.contains("show")) {
      // Handle hiding
      form.classList.remove("show");
      button.textContent = "Request a Quote";
    } else {
      // Handle showing
      form.classList.add("show");
      button.textContent = "Hide Contact Form";
    }
  });

  // Form submission handler
  window.submitForm = function (event) {
    event.preventDefault();
    const formData = new FormData(form);

    // Disable submit button to prevent double submission
    const submitButton = form.querySelector('button[type="submit"]');
    submitButton.disabled = true;
    submitButton.classList.add("opacity-50");

    fetch("/contact/post", {
      method: "POST",
      body: formData,
      headers: {
        Accept: "application/json",
      },
    })
      .then((response) => {
        if (response.redirected) {
          window.location.href = response.url;
          return;
        }

        if (response.ok) {
          // If no redirect, handle success manually
          alert("Message sent successfully!");
          form.reset();
          form.classList.remove("show");
          button.textContent = "Request a Quote";
          window.location.href = "/";
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
});
