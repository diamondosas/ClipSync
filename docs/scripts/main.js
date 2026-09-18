/**
 * ClipSync Main App Entry Point
 * Orchestrates OS detection, GitHub API, Tabs, Live LAN Sync Playground, and Navigation
 */

document.addEventListener('DOMContentLoaded', async () => {
  // 1. Initialize UI components
  if (window.UITabs) window.UITabs.init();
  if (window.ClipboardManager) window.ClipboardManager.init();

  // 2. Perform OS Detection & apply to Hero & Tabs
  if (window.OSDetector) {
    const detectedOS = window.OSDetector.detect();
    window.OSDetector.applyDetectedOS(detectedOS);
  }

  // 3. Fetch GitHub Stars & Releases
  if (window.GitHubAPI) {
    await window.GitHubAPI.init();
  }

  // 4. Initialize Live LAN Playground Simulator
  initLANPlayground();

  // 5. Mobile Menu Drawer Toggle
  const menuBtn = document.querySelector('.mobile-menu-btn');
  const drawer = document.querySelector('.mobile-drawer');
  if (menuBtn && drawer) {
    menuBtn.addEventListener('click', () => {
      drawer.classList.toggle('open');
      const isOpen = drawer.classList.contains('open');
      menuBtn.setAttribute('aria-expanded', isOpen ? 'true' : 'false');
    });

    drawer.querySelectorAll('a').forEach(link => {
      link.addEventListener('click', () => {
        drawer.classList.remove('open');
        menuBtn.setAttribute('aria-expanded', 'false');
      });
    });
  }
});

/**
 * Interactive Live LAN Sync Playground Simulator
 */
function initLANPlayground() {
  const inputEl = document.getElementById('playground-input');
  const mobilePreview = document.getElementById('playground-mobile-text');
  const speedBadge = document.getElementById('playground-speed');
  const mobileCopyBtn = document.getElementById('playground-mobile-copy');
  const sampleChips = document.querySelectorAll('.sample-chip');

  if (!inputEl || !mobilePreview) return;

  const updateSync = (text) => {
    const content = text.trim() || 'Copied text snippet will appear here instantly...';
    mobilePreview.textContent = content;

    // Simulate real-time sub-10ms LAN packet arrival
    if (speedBadge) {
      const ms = (Math.random() * 4 + 2.5).toFixed(1);
      speedBadge.textContent = `Synced in ${ms}ms`;
      speedBadge.style.color = 'var(--color-green)';
      setTimeout(() => {
        if (speedBadge) speedBadge.style.color = 'var(--color-accent-dark)';
      }, 800);
    }
  };

  inputEl.addEventListener('input', (e) => {
    updateSync(e.target.value);
  });

  sampleChips.forEach(chip => {
    chip.addEventListener('click', () => {
      const val = chip.getAttribute('data-sample') || chip.textContent;
      inputEl.value = val;
      updateSync(val);
      if (window.ClipboardManager) {
        window.ClipboardManager.showToast('Updated simulator clipboard text', 'success');
      }
    });
  });

  if (mobileCopyBtn) {
    mobileCopyBtn.addEventListener('click', () => {
      const text = mobilePreview.textContent;
      if (window.ClipboardManager && text) {
        window.ClipboardManager.copyText(text, mobileCopyBtn);
      }
    });
  }
}
