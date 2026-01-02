document.addEventListener("DOMContentLoaded", () => {
    const hero = document.getElementById('hero-slideshow');
    if (!hero) return;

    const slidesData = hero.getAttribute('data-slides');
    if (!slidesData) return;

    let images = [];
    try {
        images = JSON.parse(slidesData);
    } catch (e) {
        console.error("Invalid slides data", e);
        return;
    }

    if (images.length < 2) return;

    // Create layers for crossfading
    // Layer 1 (Active)
    const layer1 = document.createElement('div');
    layer1.className = 'absolute inset-0 bg-cover bg-center transition-opacity duration-1000 z-0';
    layer1.style.backgroundImage = `url('${images[0]}')`;
    
    // Layer 2 (Next)
    const layer2 = document.createElement('div');
    layer2.className = 'absolute inset-0 bg-cover bg-center transition-opacity duration-1000 z-0 opacity-0';
    layer2.style.backgroundImage = `url('${images[1]}')`;

    // Overlay (Gradient) - Ensure it stays on top of images but below content
    const overlay = document.createElement('div');
    overlay.className = 'absolute inset-0 z-0';
    overlay.style.backgroundImage = 'linear-gradient(to bottom, rgba(0, 0, 0, 0.2), rgba(0, 0, 0, 0.5))';

    // Content container needs z-10
    // We assume the existing content is inside the hero div. 
    // We need to insert these layers as the first children.
    hero.insertBefore(overlay, hero.firstChild);
    hero.insertBefore(layer2, hero.firstChild);
    hero.insertBefore(layer1, hero.firstChild);

    // Remove the original background image from the parent to avoid conflict
    hero.style.backgroundImage = 'none';

    let currentIndex = 0;
    let activeLayer = layer1;
    let nextLayer = layer2;

    // Preload all images
    images.forEach(src => {
        const img = new Image();
        img.src = src;
    });

    setInterval(() => {
        const nextIndex = (currentIndex + 1) % images.length;
        
        // Prepare next layer
        nextLayer.style.backgroundImage = `url('${images[nextIndex]}')`;
        
        // Fade in next, fade out active
        nextLayer.classList.remove('opacity-0');
        activeLayer.classList.add('opacity-0');

        // Swap references
        currentIndex = nextIndex;
        const temp = activeLayer;
        activeLayer = nextLayer;
        nextLayer = temp;
        
    }, 5000); // Change every 5 seconds
});
