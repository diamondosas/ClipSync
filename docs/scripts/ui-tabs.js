/**
 * ClipSync UI Tabs, Accordions, and Search Filter Module
 */

const UITabs = {
  init() {
    this.bindTabGroups();
    this.bindAccordions();
    this.bindSearchFilter();
  },

  bindTabGroups() {
    document.querySelectorAll('.tabs-group').forEach(group => {
      const tabs = group.querySelectorAll('.tab-btn');
      tabs.forEach(tab => {
        tab.addEventListener('click', () => {
          const targetKey = tab.getAttribute('data-tab');
          this.switchTabInGroup(group, targetKey);
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
      if (tab.getAttribute('data-tab') === targetKey) {
        tab.classList.add('active');
        tab.setAttribute('aria-selected', 'true');
      } else {
        tab.classList.remove('active');
        tab.setAttribute('aria-selected', 'false');
      }
    });

    panes.forEach(pane => {
      if (pane.getAttribute('data-pane') === targetKey) {
        pane.classList.add('active');
      } else {
        pane.classList.remove('active');
      }
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
    });
  }
};

window.UITabs = UITabs;
