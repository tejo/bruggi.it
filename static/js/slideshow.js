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

    // Helper to fix paths if needed (e.g. strip leading slash for relative)
    // But for now we trust the input or rely on base tag if present.
    
    // Create layers for crossfading
    const createLayer = (zIndex, opacity) => {
        const div = document.createElement('div');
        div.className = 'absolute inset-0 bg-cover bg-center';
        div.style.transition = 'opacity 1s ease-in-out';
        div.style.zIndex = zIndex;
        div.style.opacity = opacity;
        return div;
    };

    // Layer 1 (Active initially)
    const layer1 = createLayer(0, '1');
    layer1.style.backgroundImage = `url('${images[0]}')`;
    
    // Layer 2 (Next initially)
    const layer2 = createLayer(0, '0');
    layer2.style.backgroundImage = `url('${images[1]}')`;

    // Overlay (Gradient)
    const overlay = document.createElement('div');
    overlay.className = 'absolute inset-0';
    overlay.style.zIndex = '0'; // Same level, but will be behind layers due to insert order logic below?
    // Wait, we want overlay ON TOP of images but BELOW content.
    // Content is z-10. Images z-0. Overlay z-0.
    // DOM order: Images -> Overlay -> Content.
    overlay.style.backgroundImage = 'linear-gradient(to bottom, rgba(0, 0, 0, 0.2), rgba(0, 0, 0, 0.5))';

    // Insert order:
    // 1. layer1
    // 2. layer2
    // 3. overlay
    // 4. existing content (already there)
    
    // insertBefore puts new node BEFORE reference.
    // We want: [Layer1, Layer2, Overlay, Content]
    // Content is hero.firstChild.
    
    // Insert Overlay first (closest to content)
    hero.insertBefore(overlay, hero.firstChild);
    // Insert Layer 2 before Overlay
    hero.insertBefore(layer2, overlay);
    // Insert Layer 1 before Layer 2
    hero.insertBefore(layer1, layer2);

    // Remove the original background image
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
        
        // Prepare next layer (it is currently invisible)
        nextLayer.style.backgroundImage = `url('${images[nextIndex]}')`;
        
        // Swap visibility
        // nextLayer (which was behind/invisible) fades in?
        // Wait, if they are z-0, the one later in DOM (layer2) covers layer1.
        // So layer2 is ON TOP.
        
        if (nextLayer === layer2) {
            // Layer 2 is coming in. It is on top.
            // Fade it in.
            nextLayer.style.opacity = '1';
            // Layer 1 can stay 1 or go 0.
            // Better to fade out active layer?
            // If top layer fades in, it covers bottom.
            // If top layer fades out, it reveals bottom.
        } else {
            // Layer 1 is coming in. Layer 1 is BEHIND Layer 2.
            // We can't just fade Layer 1 in, because Layer 2 is opaque and covering it.
            // We must fade Layer 2 OUT to reveal Layer 1.
        }
        
        // This toggling logic is tricky with stacked layers.
        // Easier: Always fade IN the 'next' layer and make sure it's on top?
        // Or: Toggle opacity of TOP layer (layer2).
        // If Layer 2 is Opaque -> We see Image B.
        // If Layer 2 is Transparent -> We see Layer 1 (Image A).
        
        // Strategy:
        // Layer 1 always holds "Previous/Base" image.
        // Layer 2 fades in/out to show "Current/Top" image.
        
        // Let's refine the loop.
        
        if (activeLayer === layer1) {
            // We are viewing Layer 1 (Base). Layer 2 is transparent.
            // We want to transition to Image(nextIndex).
            // Put Image(nextIndex) on Layer 2.
            layer2.style.backgroundImage = `url('${images[nextIndex]}')`;
            // Fade Layer 2 IN.
            layer2.style.opacity = '1';
            
            // Now Layer 2 is Active.
            activeLayer = layer2;
            nextLayer = layer1; // Next time we will transition TO layer 1 logic
        } else {
            // We are viewing Layer 2 (Top). It is opaque.
            // We want to transition to Image(nextIndex).
            // Put Image(nextIndex) on Layer 1 (Base).
            layer1.style.backgroundImage = `url('${images[nextIndex]}')`;
            // Fade Layer 2 OUT to reveal Layer 1.
            layer2.style.opacity = '0';
            
            // Now Layer 1 is Active.
            activeLayer = layer1;
            nextLayer = layer2;
        }
        
        currentIndex = nextIndex;
        
    }, 5000);
});
