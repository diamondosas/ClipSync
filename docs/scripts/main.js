/**
 * ClipSync Main App Entry Point
 * Orchestrates OS detection, GitHub API, Tabs, and Interactive UI
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

  // 4. Mobile Menu Drawer Toggle
  const menuBtn = document.querySelector('.mobile-menu-btn');
  const drawer = document.querySelector('.mobile-drawer');
  if (menuBtn && drawer) {
    menuBtn.addEventListener('click', () => {
      drawer.classList.toggle('open');
    });

    drawer.querySelectorAll('a').forEach(link => {
      link.addEventListener('click', () => {
        drawer.classList.remove('open');
      });
    });
  }

  // 5. ScrollSpy for Active Nav Link
  const sections = document.querySelectorAll('section[id]');
  const navLinks = document.querySelectorAll('.navbar .nav-link');

  const onScroll = () => {
    const scrollY = window.pageYOffset;
    sections.forEach(current => {
      const sectionHeight = current.offsetHeight;
      const sectionTop = current.offsetTop - 100;
      const sectionId = current.getAttribute('id');

      if (scrollY > sectionTop && scrollY <= sectionTop + sectionHeight) {
        navLinks.forEach(link => {
          if (link.getAttribute('href') === `#${sectionId}`) {
            link.classList.add('active');
          } else {
            link.classList.remove('active');
          }
        });
      }
    });
  };

  window.addEventListener('scroll', onScroll);
  onScroll();
});
