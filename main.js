document.addEventListener("DOMContentLoaded", () => {
  fetch("layout.html")
    .then(response => response.text())
    .then(data => {
      const parser = new DOMParser();
      const doc = parser.parseFromString(data, "text/html");
      const header = doc.getElementById("main-header");
      const footer = doc.getElementById("main-footer");

      // Inject Header
      const headerPlaceholder = document.getElementById("layout-header");
      if (headerPlaceholder && header) {
        headerPlaceholder.replaceWith(header);
      } else if (header) {
        document.body.prepend(header);
      }

      // Inject Footer
      const footerPlaceholder = document.getElementById("layout-footer");
      if (footerPlaceholder && footer) {
        footerPlaceholder.replaceWith(footer);
      } else if (footer) {
        document.body.append(footer);
      }

      // Highlight Active Link
      let currentPath = window.location.pathname.split("/").pop();
      if (currentPath === "") currentPath = "index.html";
      
      const links = document.querySelectorAll(".nav-link");
      links.forEach(link => {
        if (link.getAttribute("href") === currentPath) {
          link.classList.add("text-primary");
        }
      });

      // Mobile Menu Logic
      const burgerBtn = document.getElementById("burger-btn");
      const mobileMenu = document.getElementById("mobile-menu");
      if (burgerBtn && mobileMenu) {
        burgerBtn.addEventListener("click", () => {
           mobileMenu.classList.toggle("hidden");
        });
      }
    })
    .catch(err => console.error("Error loading layout:", err));
});
