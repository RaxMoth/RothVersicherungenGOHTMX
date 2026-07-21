// Progressive enhancement for the fixed header: dropdown menus and the
// mobile menu toggle. Menus close on outside click and Escape.
(function () {
    var header = document.querySelector('[data-header]');
    if (!header) return;

    var dropdownButtons = header.querySelectorAll('[data-dropdown-button]');
    var mobileButton = header.querySelector('[data-mobile-button]');
    var mobileMenu = header.querySelector('#primary-menu');

    function setDropdown(btn, open) {
        btn.setAttribute('aria-expanded', String(open));
        btn.nextElementSibling.classList.toggle('hidden', !open);
        var chevron = btn.querySelector('svg');
        if (chevron) chevron.classList.toggle('rotate-180', open);
    }

    function closeDropdowns(except) {
        dropdownButtons.forEach(function (btn) {
            if (btn !== except) setDropdown(btn, false);
        });
    }

    function setMobileMenu(open) {
        if (!mobileButton || !mobileMenu) return;
        mobileButton.setAttribute('aria-expanded', String(open));
        mobileMenu.classList.toggle('hidden', !open);
        mobileButton.querySelector('[data-icon-open]').classList.toggle('hidden', open);
        mobileButton.querySelector('[data-icon-close]').classList.toggle('hidden', !open);
    }

    dropdownButtons.forEach(function (btn) {
        btn.addEventListener('click', function () {
            var open = btn.getAttribute('aria-expanded') === 'true';
            closeDropdowns(btn);
            setDropdown(btn, !open);
        });
    });

    if (mobileButton) {
        mobileButton.addEventListener('click', function () {
            setMobileMenu(mobileButton.getAttribute('aria-expanded') !== 'true');
        });
    }

    document.addEventListener('mousedown', function (e) {
        if (!header.contains(e.target)) closeDropdowns();
    });
    document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape') {
            closeDropdowns();
            setMobileMenu(false);
        }
    });
})();
