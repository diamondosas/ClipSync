/**
 * ClipSync UI Tabs, Accordions, and Category Filters
 * Accessible, keyboard-navigable tabs and search filtering
 */

const UITabs = {
  init() {
    this.bindTabGroups();
    this.bindAccordions();
    this.bindSearchFilter();
    this.bindCategoryFilterPills();
  },

  bindTabGroups() {
    document.querySelectorAll('.tabs-group').forEach(group => {
      const tabs = group.querySelectorAll('.tab-btn');
      tabs.forEach(tab => {
        tab.addEventListener('click', () => {
          const targetKey = tab.getAttribute('data-tab');
          this.switchTabInGroup(group, targetKey);
        });

        // Keyboard navigation across tabs
        tab.addEventListener('keydown', (e) => {
          const tabList = Array.from(tabs);
          const index = tabList.indexOf(tab);

          if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
            e.preventDefault();
            const nextTab = tabList[(index + 1) % tabList.length];
            nextTab.focus();
            nextTab.click();
          } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
            e.preventDefault();
            const prevTab = tabList[(index - 1 + tabList.length) % tabList.length];
            prevTab.focus();
            prevTab.click();
          }
        });
      });
    });
  },

  switchTab(groupId, tabKey) {
    const group = document.getElementById(groupId);
    if (group) {
      this.switchTabInGroup(group, tabKey);
    }
  },

  switchTabInGroup(group, targetKey) {
    const tabs = group.querySelectorAll('.tab-btn');
    const panes = group.querySelectorAll('.tab-pane');

    tabs.forEach(tab => {
      const isTarget = tab.getAttribute('data-tab') === targetKey;
      tab.classList.toggle('active', isTarget);
      tab.setAttribute('aria-selected', isTarget ? 'true' : 'false');
      tab.setAttribute('tabindex', isTarget ? '0' : '-1');
    });

    panes.forEach(pane => {
      const isTarget = pane.getAttribute('data-pane') === targetKey;
      pane.classList.toggle('active', isTarget);
    });
  },

  bindAccordions() {
    document.querySelectorAll('.accordion-header').forEach(header => {
      header.addEventListener('click', () => {
        const item = header.closest('.accordion-item');
        const isActive = item.classList.contains('active');

        // Toggle current item
        item.classList.toggle('active', !isActive);
      });
    });
  },

  bindSearchFilter() {
    const searchInput = document.getElementById('troubleshoot-search');
    if (!searchInput) return;

    searchInput.addEventListener('input', (e) => {
      const query = e.target.value.toLowerCase().trim();
      this.filterAccordionItems(query);
    });
  },

  bindCategoryFilterPills() {
    const pills = document.querySelectorAll('.filter-pill');
    if (!pills.length) return;

    pills.forEach(pill => {
      pill.addEventListener('click', () => {
        pills.forEach(p => p.classList.remove('active'));
        pill.classList.add('active');

        const category = pill.getAttribute('data-category');
        const searchInput = document.getElementById('troubleshoot-search');
        if (searchInput) searchInput.value = '';

        this.filterByCategory(category);
      });
    });
  },

  filterAccordionItems(query) {
    const items = document.querySelectorAll('.accordion-item');
    let matchCount = 0;

    items.forEach(item => {
      const title = item.querySelector('.accordion-header')?.textContent.toLowerCase() || '';
      const body = item.querySelector('.accordion-body')?.textContent.toLowerCase() || '';
      const keywords = item.getAttribute('data-keywords')?.toLowerCase() || '';

      if (!query || title.includes(query) || body.includes(query) || keywords.includes(query)) {
        item.style.display = 'block';
        matchCount++;
        if (query) {
          item.classList.add('active');
        }
      } else {
        item.style.display = 'none';
      }
    });

    const noResults = document.getElementById('no-search-results');
    if (noResults) {
      noResults.style.display = matchCount === 0 ? 'block' : 'none';
    }
  },

  filterByCategory(category) {
    const items = document.querySelectorAll('.accordion-item');
    let matchCount = 0;

    items.forEach(item => {
      const keywords = item.getAttribute('data-keywords')?.toLowerCase() || '';
      if (category === 'all' || keywords.includes(category.toLowerCase())) {
        item.style.display = 'block';
        matchCount++;
      } else {
        item.style.display = 'none';
      }
    });

    const noResults = document.getElementById('no-search-results');
    if (noResults) {
      noResults.style.display = matchCount === 0 ? 'block' : 'none';
    }
  }
};

window.UITabs = UITabs;
