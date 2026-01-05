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

    // Navigation Buttons
    const prevBtn = document.createElement('button');
    prevBtn.className = 'lightbox-nav lightbox-prev';
    prevBtn.innerHTML = '&#10094;'; // <
    prevBtn.ariaLabel = "Previous image";

    const nextBtn = document.createElement('button');
    nextBtn.className = 'lightbox-nav lightbox-next';
    nextBtn.innerHTML = '&#10095;'; // >
    nextBtn.ariaLabel = "Next image";

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
    lightbox.appendChild(prevBtn);
    lightbox.appendChild(nextBtn);
    document.body.appendChild(lightbox);
    
    let galleryItems = [];
    let currentIndex = 0;

    // Helper to get item data
    const getItemData = (element) => {
        return {
            src: element.getAttribute('href') || element.getAttribute('data-src'),
            alt: element.getAttribute('title') || element.getAttribute('data-alt'),
            author: element.getAttribute('data-author')
        };
    };

    // Update Lightbox Content
    const updateLightbox = (index) => {
        if (index < 0 || index >= galleryItems.length) return;
        currentIndex = index;
        const item = galleryItems[currentIndex];
        
        img.src = item.src;
        
        // Update Caption
        if (item.alt) {
            caption.textContent = item.alt;
            caption.style.display = 'block';
        } else {
            caption.style.display = 'none';
        }

        // Update Author
        if (item.author) {
            authorLink.href = `https://instagram.com/${item.author}`;
            authorLink.textContent = `@${item.author}`;
            authorLink.style.display = 'inline-block';
        } else {
            authorLink.style.display = 'none';
        }
    };

    // Logic to open lightbox
    const openLightbox = (index) => {
        updateLightbox(index);
        lightbox.classList.add('active');
        document.body.style.overflow = 'hidden'; // Prevent scrolling
    };
    
    // Logic to close lightbox
    const closeLightbox = () => {
        lightbox.classList.remove('active');
        document.body.style.overflow = '';
        setTimeout(() => { img.src = ''; }, 300); // Clear src after fade out
    };

    // Navigation
    const showNext = (e) => {
        e.stopPropagation();
        currentIndex = (currentIndex + 1) % galleryItems.length;
        updateLightbox(currentIndex);
    };

    const showPrev = (e) => {
        e.stopPropagation();
        currentIndex = (currentIndex - 1 + galleryItems.length) % galleryItems.length;
        updateLightbox(currentIndex);
    };
    
    // Event Listeners
    closeBtn.addEventListener('click', closeLightbox);
    nextBtn.addEventListener('click', showNext);
    prevBtn.addEventListener('click', showPrev);
    
    lightbox.addEventListener('click', (e) => {
        if (e.target === lightbox) {
            closeLightbox();
        }
    });
    
    document.addEventListener('keydown', (e) => {
        if (!lightbox.classList.contains('active')) return;
        
        if (e.key === 'Escape') closeLightbox();
        if (e.key === 'ArrowRight') showNext(e);
        if (e.key === 'ArrowLeft') showPrev(e);
    });

    // Touch Swipe Support
    let touchStartX = 0;
    let touchEndX = 0;

    lightbox.addEventListener('touchstart', (e) => {
        touchStartX = e.changedTouches[0].screenX;
    }, { passive: true });

    lightbox.addEventListener('touchend', (e) => {
        touchEndX = e.changedTouches[0].screenX;
        handleSwipe();
    }, { passive: true });

    const handleSwipe = () => {
        const swipeThreshold = 50; // Minimum distance for swipe
        if (touchEndX < touchStartX - swipeThreshold) {
            showNext({ stopPropagation: () => {} }); // Swipe Left -> Next
        }
        if (touchEndX > touchStartX + swipeThreshold) {
            showPrev({ stopPropagation: () => {} }); // Swipe Right -> Prev
        }
    };
    
    // Attach to images with 'lightbox-trigger' class
    // We re-scan on click to support dynamic content if needed, but for now scan once
    const triggers = Array.from(document.querySelectorAll('.lightbox-trigger'));
    
    // Populate gallery items list
    galleryItems = triggers.map(getItemData);

    triggers.forEach((trigger, index) => {
        trigger.addEventListener('click', (e) => {
            e.preventDefault();
            // Re-fetch items in case of dynamic updates (optional but safe)
            // galleryItems = Array.from(document.querySelectorAll('.lightbox-trigger')).map(getItemData);
            openLightbox(index);
        });
    });
});
