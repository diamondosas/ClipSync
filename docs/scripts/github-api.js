/**
 * ClipSync GitHub API Integration Module
 * Fetches real-time stars and latest release artifacts from GitHub Actions/Releases
 */

const GitHubAPI = {
  repoOwner: 'DiamondOsas',
  repoName: 'ClipSync',
  defaultVersion: 'v1.0.0',

  get baseUrl() {
    return `https://api.github.com/repos/${this.repoOwner}/${this.repoName}`;
  },

  async init() {
    await Promise.allSettled([
      this.fetchRepoStats(),
      this.fetchLatestRelease()
    ]);
  },

  async fetchRepoStats() {
    try {
      const response = await fetch(this.baseUrl, {
        headers: { 'Accept': 'application/vnd.github.v3+json' }
      });
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      
      const data = await response.json();
      const stars = data.stargazers_count ?? 0;
      this.updateStarElements(stars);
    } catch (err) {
      console.warn('[ClipSync] GitHub stats API fetch error (using fallback):', err.message);
      this.updateStarElements(null);
    }
  },

  updateStarElements(count) {
    const starEls = document.querySelectorAll('.star-count');
    starEls.forEach(el => {
      if (count !== null && count !== undefined) {
        el.textContent = count >= 1000 ? (count / 1000).toFixed(1) + 'k' : count.toString();
      } else {
        el.textContent = 'Star';
      }
    });
  },

  async fetchLatestRelease() {
    try {
      const response = await fetch(`${this.baseUrl}/releases/latest`, {
        headers: { 'Accept': 'application/vnd.github.v3+json' }
      });
      if (!response.ok) throw new Error(`HTTP ${response.status}`);

      const release = await response.json();
      this.applyReleaseData(release);
    } catch (err) {
      console.warn('[ClipSync] GitHub releases API fetch error (using fallback):', err.message);
      this.applyFallbackReleaseData();
    }
  },

  applyReleaseData(release) {
    const tagName = release.tag_name || this.defaultVersion;
    const publishedAt = release.published_at ? new Date(release.published_at).toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }) : 'Latest Release';

    // Update version badges across the UI
    document.querySelectorAll('.release-version-tag').forEach(el => {
      el.textContent = tagName;
    });

    document.querySelectorAll('.release-date-tag').forEach(el => {
      el.textContent = publishedAt;
    });

    // Map release assets to download buttons
    const assetsMap = {};
    if (Array.isArray(release.assets)) {
      release.assets.forEach(asset => {
        assetsMap[asset.name] = {
          url: asset.browser_download_url,
          size: this.formatBytes(asset.size)
        };
      });
    }

    // Update buttons with direct asset URLs
    document.querySelectorAll('[data-target-file]').forEach(btn => {
      const fileName = btn.getAttribute('data-target-file');
      if (assetsMap[fileName]) {
        btn.setAttribute('href', assetsMap[fileName].url);
        const sizeBadge = btn.querySelector('.asset-size-badge');
        if (sizeBadge) {
          sizeBadge.textContent = assetsMap[fileName].size;
        }
      }
    });
  },

  applyFallbackReleaseData() {
    const defaultTag = this.defaultVersion;
    document.querySelectorAll('.release-version-tag').forEach(el => {
      el.textContent = defaultTag;
    });

    document.querySelectorAll('[data-target-file]').forEach(btn => {
      const fileName = btn.getAttribute('data-target-file');
      const directUrl = `https://github.com/${this.repoOwner}/${this.repoName}/releases/latest/download/${fileName}`;
      btn.setAttribute('href', directUrl);
    });
  },

  formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }
};

window.GitHubAPI = GitHubAPI;
