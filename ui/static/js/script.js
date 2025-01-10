document
  .getElementById("showContactForm")
  .addEventListener("click", function () {
    const contactForm = document.getElementById("contactForm");
    contactForm.classList.toggle("show");
    this.textContent = contactForm.classList.contains("show")
      ? "Hide Contact Form"
      : "Request a Quote";
  });
