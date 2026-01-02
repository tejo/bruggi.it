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

    // Info Container (Caption + Author)
    const infoContainer = document.createElement('div');
    infoContainer.className = 'lightbox-info';
    
    const caption = document.createElement('p');
    caption.className = 'lightbox-caption';
    
    const authorLink = document.createElement('a');
    authorLink.className = 'lightbox-author';
    authorLink.target = '_blank';
    authorLink.rel = 'noopener noreferrer';
    
    infoContainer.appendChild(caption);
    infoContainer.appendChild(authorLink);
    
    lightbox.appendChild(img);
    lightbox.appendChild(infoContainer);
    lightbox.appendChild(closeBtn);
    document.body.appendChild(lightbox);
    
    // Logic to open lightbox
    const openLightbox = (src, altText, authorHandle) => {
        img.src = src;
        lightbox.classList.add('active');
        document.body.style.overflow = 'hidden'; // Prevent scrolling
        
        // Update Caption
        if (altText) {
            caption.textContent = altText;
            caption.style.display = 'block';
        } else {
            caption.style.display = 'none';
        }

        // Update Author
        if (authorHandle) {
            authorLink.href = `https://instagram.com/${authorHandle}`;
            authorLink.textContent = `@${authorHandle}`;
            authorLink.style.display = 'inline-block';
        } else {
            authorLink.style.display = 'none';
        }
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
    const triggers = document.querySelectorAll('.lightbox-trigger');
    triggers.forEach(trigger => {
        trigger.addEventListener('click', (e) => {
            e.preventDefault();
            const src = trigger.getAttribute('href') || trigger.getAttribute('data-src');
            const alt = trigger.getAttribute('title') || trigger.getAttribute('data-alt');
            const author = trigger.getAttribute('data-author');
            
            if (src) openLightbox(src, alt, author);
        });
    });
});
