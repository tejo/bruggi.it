document.addEventListener("DOMContentLoaded", () => {
    // Create Lightbox DOM
    const lightbox = document.createElement('div');
    lightbox.className = 'lightbox-modal';
    
    const img = document.createElement('img');
    img.className = 'lightbox-content';
    
    const closeBtn = document.createElement('button');
    closeBtn.className = 'lightbox-close';
    closeBtn.innerHTML = '&times;';
    closeBtn.ariaLabel = "Close lightbox";
    
    lightbox.appendChild(img);
    lightbox.appendChild(closeBtn);
    document.body.appendChild(lightbox);
    
    // Logic to open lightbox
    const openLightbox = (src) => {
        img.src = src;
        lightbox.classList.add('active');
        document.body.style.overflow = 'hidden'; // Prevent scrolling
    };
    
    // Logic to close lightbox
    const closeLightbox = () => {
        lightbox.classList.remove('active');
        document.body.style.overflow = '';
        setTimeout(() => { img.src = ''; }, 300); // Clear src after fade out
    };
    
    // Event Listeners
    closeBtn.addEventListener('click', closeLightbox);
    
    lightbox.addEventListener('click', (e) => {
        if (e.target === lightbox) {
            closeLightbox();
        }
    });
    
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && lightbox.classList.contains('active')) {
            closeLightbox();
        }
    });
    
    // Attach to images with 'lightbox-trigger' class
    // We delegate the event to the body to handle dynamically loaded content if any (though this is static site)
    // Or just attach to existing ones.
    const triggers = document.querySelectorAll('.lightbox-trigger');
    triggers.forEach(trigger => {
        trigger.addEventListener('click', (e) => {
            e.preventDefault();
            const src = trigger.getAttribute('href') || trigger.getAttribute('data-src');
            if (src) openLightbox(src);
        });
    });
});
