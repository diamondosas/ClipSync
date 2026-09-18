/**
 * ClipSync Clipboard & Toast Utilities
 * Enables 1-click code copying, animated toast notifications, and debug info generator
 */

const ClipboardManager = {
  init() {
    this.createToastContainer();
    this.bindCopyButtons();
    this.bindDebugInfoGenerator();
  },

  createToastContainer() {
    if (!document.querySelector('.toast-container')) {
      const container = document.createElement('div');
      container.className = 'toast-container';
      document.body.appendChild(container);
    }
  },

  bindCopyButtons() {
    document.addEventListener('click', (e) => {
      const btn = e.target.closest('.btn-copy');
      if (!btn) return;

      e.preventDefault();
      const targetId = btn.getAttribute('data-copy-target');
      let textToCopy = '';

      if (targetId) {
        const targetEl = document.getElementById(targetId);
        if (targetEl) {
          textToCopy = targetEl.innerText.replace(/^\$\s*/gm, '').replace(/^PS>\s*/gm, '').trim();
        }
      } else if (btn.getAttribute('data-copy-text')) {
        textToCopy = btn.getAttribute('data-copy-text').trim();
      } else {
        const codeEl = btn.closest('.code-block')?.querySelector('code');
        if (codeEl) {
          textToCopy = codeEl.innerText.replace(/^\$\s*/gm, '').replace(/^PS>\s*/gm, '').trim();
        }
      }

      if (textToCopy) {
        this.copyText(textToCopy, btn);
      }
    });
  },

  bindDebugInfoGenerator() {
    document.addEventListener('click', (e) => {
      const debugBtn = e.target.closest('#btn-copy-debug-info');
      if (!debugBtn) return;

      e.preventDefault();
      const detectedOS = window.OSDetector ? window.OSDetector.detect() : { name: 'Unknown' };
      const debugReport = `
=== ClipSync Debug Information ===
Timestamp: ${new Date().toISOString()}
Detected Platform: ${detectedOS.name}
User Agent: ${navigator.userAgent}
Screen Resolution: ${window.screen.width}x${window.screen.height}
Default Ports: UDP 9999 (KCP), UDP 5353 (mDNS Discovery)
Issue Summary: 
==================================
`.trim();

      this.copyText(debugReport, debugBtn);
    });
  },

  async copyText(text, triggerBtn) {
    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(text);
      } else {
        // Fallback for non-https / older contexts
        const textarea = document.createElement('textarea');
        textarea.value = text;
        textarea.style.position = 'fixed';
        textarea.style.left = '-9999px';
        document.body.appendChild(textarea);
        textarea.focus();
        textarea.select();
        document.execCommand('copy');
        textarea.remove();
      }
      
      if (triggerBtn) {
        const originalHtml = triggerBtn.innerHTML;
        triggerBtn.classList.add('copied');
        triggerBtn.innerHTML = `
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span>Copied!</span>
        `;

        setTimeout(() => {
          triggerBtn.classList.remove('copied');
          triggerBtn.innerHTML = originalHtml;
        }, 2200);
      }

      this.showToast('Copied to clipboard!', 'success');
    } catch (err) {
      console.error('Failed to copy text:', err);
      this.showToast('Failed to copy. Please select and copy manually.', 'error');
    }
  },

  showToast(message, type = 'success') {
    const container = document.querySelector('.toast-container');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast ${type === 'error' ? 'toast-error' : ''}`;
    
    const icon = type === 'success' 
      ? `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#15803D" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>`
      : `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#C53030" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`;

    toast.innerHTML = `${icon}<span>${message}</span>`;
    container.appendChild(toast);

    setTimeout(() => {
      toast.style.opacity = '0';
      toast.style.transform = 'translateY(12px) scale(0.95)';
      toast.style.transition = 'all 0.25s var(--ease-spring)';
      setTimeout(() => toast.remove(), 250);
    }, 2800);
  }
};

window.ClipboardManager = ClipboardManager;
