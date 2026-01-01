document.addEventListener("DOMContentLoaded", () => {
  // Highlight Active Link
  let currentPath = window.location.pathname.split("/").pop();
  if (currentPath === "") currentPath = "index.html";
  
  const links = document.querySelectorAll(".nav-link");
  links.forEach(link => {
    // Check if href matches current path (considering absolute/relative)
    const href = link.getAttribute("href");
    if (href && (href === currentPath || href.endsWith("/" + currentPath))) {
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
});