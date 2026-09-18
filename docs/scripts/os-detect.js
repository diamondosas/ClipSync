/**
 * ClipSync OS Auto-Detection Module
 * Detects user platform (Windows, macOS, Linux, Android) and architecture
 */

const OSDetector = {
  detect() {
    const userAgent = window.navigator.userAgent.toLowerCase();
    const platform = (window.navigator.platform || '').toLowerCase();
    const userAgentData = window.navigator.userAgentData;

    // Check Android
    if (/android/i.test(userAgent)) {
      return {
        key: 'android',
        name: 'Android',
        badge: 'Android 8.0+ (ARM64 / APK)',
        primaryExt: '.apk',
        primaryFile: 'ClipSync.apk',
        primaryLabel: 'Download for Android',
        subtext: 'Direct APK installation for mobile',
        icon: 'android'
      };
    }

    // Check macOS
    if (/macintosh|mac os x|mac_powerpc/i.test(userAgent) || /mac/i.test(platform)) {
      const isAppleSilicon = /macintosh/i.test(userAgent) && (navigator.maxTouchPoints > 0 || /arm/i.test(userAgent));
      return {
        key: 'macos',
        name: isAppleSilicon ? 'macOS (Apple Silicon)' : 'macOS (Universal)',
        badge: 'macOS 11.0+ (Universal Binary)',
        primaryExt: '.tar.gz',
        primaryFile: 'clipsync-macos.tar.gz',
        primaryLabel: 'Download for macOS',
        subtext: 'Universal binary (Apple Silicon & Intel)',
        icon: 'apple'
      };
    }

    // Check Linux
    if (/linux/i.test(userAgent) || /linux/i.test(platform)) {
      return {
        key: 'linux',
        name: 'Linux',
        badge: 'Linux (Debian / Ubuntu / Arch / Fedora x64)',
        primaryExt: '.deb',
        primaryFile: 'clipsync-linux.deb',
        primaryLabel: 'Download for Linux (.deb)',
        secondaryFile: 'clipsync-linux.tar.gz',
        secondaryLabel: 'Download Linux Tarball (.tar.gz)',
        subtext: 'For Debian, Ubuntu, Mint & generic distros',
        icon: 'linux'
      };
    }

    // Check Windows (or default)
    const is64 = /win64|x64|wow64/i.test(userAgent) || (userAgentData && userAgentData.architecture === 'x86');
    return {
      key: 'windows',
      name: 'Windows',
      badge: is64 ? 'Windows 10 / 11 (64-bit)' : 'Windows (x86 / x64)',
      primaryExt: '.exe',
      primaryFile: 'ClipSync-Setup.exe',
      primaryLabel: 'Download for Windows',
      secondaryFile: 'clipsync-windows.zip',
      secondaryLabel: 'Download Portable Zip (.zip)',
      subtext: 'Standard Inno Setup Installer (.exe)',
      icon: 'windows'
    };
  },

  applyDetectedOS(osInfo) {
    const badgeEl = document.getElementById('detected-os-text');
    const heroBtn = document.getElementById('hero-download-btn');
    const heroBtnText = document.getElementById('hero-btn-title');
    const heroBtnSub = document.getElementById('hero-btn-sub');
    const heroSecondaryBtn = document.getElementById('hero-secondary-btn');

    if (badgeEl) {
      badgeEl.textContent = osInfo.badge;
    }

    if (heroBtn && heroBtnText) {
      heroBtnText.textContent = osInfo.primaryLabel;
      if (heroBtnSub) heroBtnSub.textContent = osInfo.subtext;
      
      const downloadUrl = `https://github.com/DiamondOsas/ClipSync/releases/latest/download/${osInfo.primaryFile}`;
      heroBtn.setAttribute('href', downloadUrl);
      heroBtn.setAttribute('data-target-file', osInfo.primaryFile);
    }

    if (heroSecondaryBtn) {
      if (osInfo.secondaryFile) {
        heroSecondaryBtn.style.display = 'inline-flex';
        heroSecondaryBtn.textContent = osInfo.secondaryLabel;
        heroSecondaryBtn.setAttribute('href', `https://github.com/DiamondOsas/ClipSync/releases/latest/download/${osInfo.secondaryFile}`);
        heroSecondaryBtn.setAttribute('data-target-file', osInfo.secondaryFile);
      } else {
        heroSecondaryBtn.style.display = 'none';
      }
    }

    // Switch OS Matrix tab to the detected OS by default, leaving others visible and clickable
    if (window.UITabs && typeof window.UITabs.switchTab === 'function') {
      window.UITabs.switchTab('os-tabs', osInfo.key);
    }
  }
};

window.OSDetector = OSDetector;
