/* ============================================
   Main JavaScript - Interactive Animations
   ============================================ */

// Particle System
class ParticleSystem {
  constructor() {
    this.particles = [];
    this.particleCount = 50;
    this.init();
  }

  init() {
    const container = document.createElement('div');
    container.className = 'particles';
    document.body.appendChild(container);

    for (let i = 0; i < this.particleCount; i++) {
      const particle = document.createElement('div');
      particle.className = 'particle';
      particle.style.left = Math.random() * 100 + '%';
      particle.style.animationDelay = Math.random() * 20 + 's';
      particle.style.animationDuration = (15 + Math.random() * 10) + 's';
      container.appendChild(particle);
      this.particles.push(particle);
    }
  }
}

// Scroll Animations
class ScrollAnimations {
  constructor() {
    this.elements = document.querySelectorAll('.animate-on-scroll');
    this.init();
  }

  init() {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
        }
      });
    }, {
      threshold: 0.1,
      rootMargin: '50px'
    });

    this.elements.forEach(el => observer.observe(el));
  }
}

// Header Scroll Effect
class HeaderScroll {
  constructor() {
    this.header = document.querySelector('.header');
    this.lastScroll = 0;
    this.init();
  }

  init() {
    window.addEventListener('scroll', () => {
      const currentScroll = window.pageYOffset;

      if (currentScroll > 100) {
        this.header.classList.add('scrolled');
      } else {
        this.header.classList.remove('scrolled');
      }

      this.lastScroll = currentScroll;
    });
  }
}

// Parallax Effect
class ParallaxEffect {
  constructor() {
    this.elements = document.querySelectorAll('.parallax');
    this.init();
  }

  init() {
    window.addEventListener('scroll', () => {
      const scrolled = window.pageYOffset;

      this.elements.forEach(el => {
        const speed = el.dataset.speed || 0.5;
        const yPos = -(scrolled * speed);
        el.style.transform = `translateY(${yPos}px)`;
      });
    });
  }
}

// Mobile Navigation
class MobileNav {
  constructor() {
    this.toggle = document.querySelector('.nav-toggle');
    this.menu = document.querySelector('.nav-menu');
    this.links = document.querySelectorAll('.nav-link');
    this.init();
  }

  init() {
    if (!this.toggle || !this.menu) return;

    this.toggle.addEventListener('click', () => {
      this.toggle.classList.toggle('active');
      this.menu.classList.toggle('active');
    });

    this.links.forEach(link => {
      link.addEventListener('click', () => {
        this.toggle.classList.remove('active');
        this.menu.classList.remove('active');
      });
    });
  }
}

// Magnetic Hover Effect
class MagneticHover {
  constructor() {
    this.elements = document.querySelectorAll('.card, .card-code');
    this.init();
  }

  init() {
    this.elements.forEach(el => {
      el.addEventListener('mousemove', (e) => {
        const rect = el.getBoundingClientRect();
        const x = e.clientX - rect.left - rect.width / 2;
        const y = e.clientY - rect.top - rect.height / 2;

        const angle = Math.atan2(y, x) * (180 / Math.PI);
        const distance = Math.sqrt(x * x + y * y);
        const maxDistance = Math.max(rect.width, rect.height) / 2;
        const strength = Math.min(distance / maxDistance, 1);

        el.style.transform = `
          perspective(1000px)
          rotateX(${y * 0.05}deg)
          rotateY(${x * 0.05}deg)
          translateZ(${strength * 20}px)
        `;
      });

      el.addEventListener('mouseleave', () => {
        el.style.transform = '';
      });
    });
  }
}

// Smooth Scroll
class SmoothScroll {
  constructor() {
    this.links = document.querySelectorAll('a[href^="#"]');
    this.init();
  }

  init() {
    this.links.forEach(link => {
      link.addEventListener('click', (e) => {
        const href = link.getAttribute('href');
        if (href === '#') return;

        e.preventDefault();

        // Get the ID without the #
        const id = href.substring(1);

        // Try multiple methods to find the target element
        let target = null;

        // Method 1: Try getElementById (most reliable)
        target = document.getElementById(id);

        // Method 2: Try querySelector with various selectors
        if (!target) {
          try {
            // Try direct selector first
            target = document.querySelector(`#${CSS.escape(id)}`);
          } catch (error) {
            // Ignore CSS.escape errors
          }
        }

        // Method 3: Try variations of the ID (common Eleventy transformations)
        if (!target) {
          const variations = [
            id.toLowerCase(),
            id.replace(/\./g, ''),
            id.replace(/\./g, '-'),
            id.replace(/[^a-zA-Z0-9-_]/g, '-'),
            id.replace(/[^a-zA-Z0-9-_]/g, '').toLowerCase(),
            id.replace(/^\d+\.?-?/, ''), // Remove leading numbers and dots
            id.split('.').pop(), // Get the part after the last dot
            id.replace(/\d+\./g, ''), // Remove all number-dot patterns
            id.replace(/\./g, '').replace(/-/g, ''), // Remove both dots and dashes
            id.replace(/^(\d+)\.(\d*)-?/, '$1$2-'), // Convert 2.4-operators to 24-operators
            id.replace(/(\d+)\.(\d+)-(.+)/, '$1$2-$3'), // 2.4-operators -> 24-operators
          ];

          for (const variant of variations) {
            if (variant) {
              target = document.getElementById(variant);
              if (target) break;

              try {
                target = document.querySelector(`#${CSS.escape(variant)}`);
                if (target) break;
              } catch (error) {
                // Continue to next variation
              }
            }
          }
        }

        // Method 4: Find by heading text content
        if (!target) {
          const headingText = id.replace(/^\d+\.?\d*-?/, '').replace(/-/g, ' ');
          const headers = document.querySelectorAll('h1, h2, h3, h4, h5, h6');

          for (const header of headers) {
            const headerText = header.textContent.toLowerCase().trim();
            if (headerText.includes(headingText.toLowerCase()) ||
              headerText.replace(/[^a-zA-Z0-9\s]/g, '').includes(headingText.toLowerCase())) {
              target = header;
              break;
            }
          }
        }

        // Method 5: Find by data attributes or class names
        if (!target) {
          target = document.querySelector(`[data-anchor="${id}"]`) ||
            document.querySelector(`[data-id="${id}"]`) ||
            document.querySelector(`.${id.replace(/[^a-zA-Z0-9-_]/g, '-')}`);
        }

        if (target) {
          const offset = 100;
          const targetPosition = target.offsetTop - offset;

          window.scrollTo({
            top: targetPosition,
            behavior: 'smooth'
          });

          // Update URL hash
          history.pushState(null, null, href);
          console.log('✅ Found target for:', href, '-> Element:', target.tagName, target.id || target.className);
        } else {
          // Debug: Log all available IDs on the page
          const allIds = Array.from(document.querySelectorAll('[id]')).map(el => el.id);
          console.warn('❌ Could not find target element for:', href, 'ID:', id);
          console.warn('🔍 Available IDs on page:', allIds.filter(id => id.includes('operator') || id.includes('lexical')));

          // Fallback: just update the URL and let the browser handle it
          window.location.hash = href;
        }
      });
    });
  }
}

// Typing Animation for Code Blocks with Syntax Highlighting
class TypeWriter {
  constructor() {
    // Only apply typewriter effect to code blocks within elements that have the 'typewriter-enabled' class
    this.codeBlocks = document.querySelectorAll('.typewriter-enabled pre code.language-osprey');
    this.init();
  }

  init() {
    this.codeBlocks.forEach((block, blockIndex) => {
      // The HTML already has syntax highlighting spans!
      // Don't destroy them, just animate them appearing
      this.prepareTypingAnimation(block, blockIndex);
    });
  }

  prepareTypingAnimation(element, blockIndex) {
    // Get all text nodes and wrap them in spans for animation
    this.wrapTextNodes(element);

    // Get all character spans
    const chars = element.querySelectorAll('.typing-char');

    // Hide all characters initially
    chars.forEach(char => {
      char.style.opacity = '0';
    });

    // Start typing animation
    setTimeout(() => {
      this.animateTyping(chars);
    }, 1000 + (blockIndex * 200));
  }

  wrapTextNodes(element) {
    const walker = document.createTreeWalker(
      element,
      NodeFilter.SHOW_TEXT,
      null,
      false
    );

    const textNodes = [];
    let node;
    while (node = walker.nextNode()) {
      if (node.textContent.trim().length > 0) {
        textNodes.push(node);
      }
    }

    textNodes.forEach(textNode => {
      const parent = textNode.parentNode;
      const text = textNode.textContent;
      const fragment = document.createDocumentFragment();

      for (let i = 0; i < text.length; i++) {
        const char = text.charAt(i);
        const span = document.createElement('span');
        span.className = 'typing-char';
        span.textContent = char;
        fragment.appendChild(span);
      }

      parent.replaceChild(fragment, textNode);
    });
  }

  animateTyping(chars) {
    chars.forEach((char, index) => {
      setTimeout(() => {
        char.style.opacity = '1';
        char.style.transition = 'opacity 0.1s ease';
      }, index * 25);
    });
  }
}

// Glitch Effect on Hover
class GlitchEffect {
  constructor() {
    this.elements = document.querySelectorAll('.hero-title');
    this.init();
  }

  init() {
    this.elements.forEach(el => {
      el.addEventListener('mouseenter', () => {
        el.classList.add('glitch');
        el.setAttribute('data-text', el.textContent);
      });

      el.addEventListener('mouseleave', () => {
        setTimeout(() => {
          el.classList.remove('glitch');
        }, 300);
      });
    });
  }
}

// Cursor Glow Effect
class CursorGlow {
  constructor() {
    this.glow = document.createElement('div');
    this.init();
  }

  init() {
    this.glow.className = 'cursor-glow';
    document.body.appendChild(this.glow);

    document.addEventListener('mousemove', (e) => {
      this.glow.style.left = e.clientX + 'px';
      this.glow.style.top = e.clientY + 'px';
    });

    // Add CSS for cursor glow
    const style = document.createElement('style');
    style.textContent = `
      .cursor-glow {
        position: fixed;
        width: 400px;
        height: 400px;
        background: radial-gradient(circle, rgba(102, 126, 234, 0.1) 0%, transparent 70%);
        border-radius: 50%;
        pointer-events: none;
        transform: translate(-50%, -50%);
        z-index: -1;
        transition: opacity 0.3s ease;
        mix-blend-mode: screen;
      }
    `;
    document.head.appendChild(style);
  }
}

// Performance Optimization
function throttle(func, limit) {
  let inThrottle;
  return function () {
    const args = arguments;
    const context = this;
    if (!inThrottle) {
      func.apply(context, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}

// Initialize Everything
document.addEventListener('DOMContentLoaded', () => {
  // Create particle system
  new ParticleSystem();

  // Initialize animations
  new ScrollAnimations();
  new HeaderScroll();
  new ParallaxEffect();
  new MobileNav();
  new MagneticHover();
  new SmoothScroll();
  new GlitchEffect();
  new CursorGlow();

  // Initialize typing animation with delay
  setTimeout(() => {
    new TypeWriter();
  }, 500);

  // Add loaded class for animations
  document.body.classList.add('loaded');
});

// Preloader
window.addEventListener('load', () => {
  const preloader = document.querySelector('.preloader');
  if (preloader) {
    preloader.classList.add('fade-out');
    setTimeout(() => {
      preloader.style.display = 'none';
    }, 500);
  }
});

// Intersection Observer for lazy loading
const lazyImages = document.querySelectorAll('img[data-src]');
const imageObserver = new IntersectionObserver((entries, observer) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      const img = entry.target;
      img.src = img.dataset.src;
      img.classList.add('loaded');
      observer.unobserve(img);
    }
  });
});

lazyImages.forEach(img => imageObserver.observe(img));

// Add some console easter eggs
console.log('%c🦅 Welcome to Osprey!', 'font-size: 24px; font-weight: bold; color: #667eea;');
console.log('%cBuilding the future of programming, one function at a time.', 'font-size: 14px; color: #764ba2;');
console.log('%cInterested in contributing? Check out https://github.com/MelbourneDeveloper/osprey', 'font-size: 12px; color: #666;'); 