// Color Switcher — applies accent color overrides from localStorage.
// Persists across page loads by setting CSS custom properties on :root.

(function () {
  'use strict';

  var DEFAULT_ACCENT = '#f59e0b';
  var DEFAULT_ACCENT_LIGHT = '#fbbf24';

  function applyColors(accent, accentLight) {
    document.documentElement.style.setProperty('--color-accent', accent);
    document.documentElement.style.setProperty('--color-accent-light', accentLight);
  }

  function markActive(accent) {
    var swatches = document.querySelectorAll('[data-accent]');
    swatches.forEach(function (el) {
      var ring = el.querySelector('.color-swatch-ring');
      if (!ring) return;
      if (el.getAttribute('data-accent') === accent) {
        ring.classList.add('border-text');
        ring.classList.remove('border-bg-border');
      } else {
        ring.classList.remove('border-text');
        ring.classList.add('border-bg-border');
      }
    });
  }

  // Apply saved colors on every page load.
  var savedAccent = localStorage.getItem('accent');
  var savedAccentLight = localStorage.getItem('accent-light');
  if (savedAccent && savedAccentLight) {
    applyColors(savedAccent, savedAccentLight);
  }

  // Once DOM is ready, wire up clicks and mark active states.
  document.addEventListener('DOMContentLoaded', function () {
    var current = localStorage.getItem('accent') || DEFAULT_ACCENT;
    markActive(current);

    document.querySelectorAll('[data-accent]').forEach(function (el) {
      el.addEventListener('click', function () {
        var accent = el.getAttribute('data-accent');
        var accentLight = el.getAttribute('data-accent-light');
        applyColors(accent, accentLight);
        localStorage.setItem('accent', accent);
        localStorage.setItem('accent-light', accentLight);
        markActive(accent);
      });
    });

    var resetBtn = document.getElementById('color-reset');
    if (resetBtn) {
      resetBtn.addEventListener('click', function () {
        localStorage.removeItem('accent');
        localStorage.removeItem('accent-light');
        applyColors(DEFAULT_ACCENT, DEFAULT_ACCENT_LIGHT);
        markActive(DEFAULT_ACCENT);
      });
    }
  });
})();
