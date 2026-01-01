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

  // Itinerary Filtering
  const filterButtons = document.querySelectorAll(".filter-btn");
  const itineraryCards = document.querySelectorAll(".itinerary-card");

  if (filterButtons.length > 0) {
    filterButtons.forEach(btn => {
      btn.addEventListener("click", () => {
        const filter = btn.getAttribute("data-filter");

        // Update active button style
        filterButtons.forEach(b => {
          b.classList.remove("active", "bg-[#111811]", "dark:bg-white", "text-white", "dark:text-[#111811]");
          b.classList.add("bg-white", "dark:bg-[#2a402a]", "text-[#111811]", "dark:text-gray-200");
          b.classList.replace("font-bold", "font-medium");
        });

        btn.classList.add("active", "bg-[#111811]", "dark:bg-white", "text-white", "dark:text-[#111811]");
        btn.classList.remove("bg-white", "dark:bg-[#2a402a]", "text-[#111811]", "dark:text-gray-200");
        btn.classList.replace("font-medium", "font-bold");

        // Filter cards
        itineraryCards.forEach(card => {
          if (filter === "all" || card.getAttribute("data-type") === filter) {
            card.style.display = "flex";
          } else {
            card.style.display = "none";
          }
        });
      });
    });
  }
});