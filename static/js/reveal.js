// Reveal-on-scroll: fade + lift. JS toggles `.is-in` on `.reveal` elements
// as they enter the viewport. Loaded site-wide via base layout — pages
// without `.reveal` elements no-op (the observer just sees zero targets).
(function () {
	function activate() {
		var els = document.querySelectorAll('.reveal');
		if (!('IntersectionObserver' in window)) {
			els.forEach(function (el) { el.classList.add('is-in'); });
			return;
		}
		var io = new IntersectionObserver(function (entries) {
			entries.forEach(function (e) {
				if (e.isIntersecting) {
					e.target.classList.add('is-in');
					io.unobserve(e.target);
				}
			});
		}, { rootMargin: '0px 0px -10% 0px', threshold: 0.05 });
		els.forEach(function (el) { io.observe(el); });
	}

	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', activate);
	} else {
		activate();
	}
})();
