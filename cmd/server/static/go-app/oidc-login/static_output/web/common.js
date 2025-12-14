// Common JavaScript utilities for the application
// Log build version information when this script loads
(function() {
    if (window.BUILD_VERSION) {
        console.log('[BUILD INFO] Version:', window.BUILD_VERSION.version);
        console.log('[BUILD INFO] Build Time:', window.BUILD_VERSION.buildTime);
        console.log('[BUILD INFO] Git Commit:', window.BUILD_VERSION.gitCommit);
        console.log('[BUILD INFO] Git Branch:', window.BUILD_VERSION.gitBranch);
    } else {
        console.log('[BUILD INFO] Build version information not available yet');
    }
    
    // Example: Reading build version from meta tags (alternative to window.BUILD_VERSION)
    const versionMeta = document.querySelector('meta[name="app-version"]');
    const buildTimeMeta = document.querySelector('meta[name="app-build-time"]');
    const gitCommitMeta = document.querySelector('meta[name="app-git-commit"]');
    const gitBranchMeta = document.querySelector('meta[name="app-git-branch"]');
    
    if (versionMeta) {
        console.log('[BUILD INFO from META] Version:', versionMeta.content);
        console.log('[BUILD INFO from META] Build Time:', buildTimeMeta?.content);
        console.log('[BUILD INFO from META] Git Commit:', gitCommitMeta?.content);
        console.log('[BUILD INFO from META] Git Branch:', gitBranchMeta?.content);
        
        // Example: Use version for cache-busting in dynamic resource loading
        const version = versionMeta.content || Date.now();
        console.log('[BUILD INFO] Using version for cache-busting:', version);
        
        // Example: Load a dynamic script with cache-busting
        // const script = document.createElement('script');
        // script.src = `/my-dynamic-script.js?v=${version}`;
        // document.head.appendChild(script);
    }
})();

// Global config storage
window.appConfig = null;
window.appConfigLoaded = false;
window.appConfigError = null;

// Fetch app config before WASM loads
(async function fetchAppConfig() {
  try {
    console.log("Fetching app config...");
    // Use build version if available, otherwise use timestamp
    const version = window.BUILD_VERSION?.version || Date.now();
    const response = await fetch(`web/app.json?v=${version}`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch config: ${response.status} ${response.statusText}`);
    }
    
    window.appConfig = await response.json();
    window.appConfigLoaded = true;
    console.log("App config loaded successfully:", window.appConfig);
  } catch (error) {
    console.error("Failed to load app config:", error);
    window.appConfigError = error.message;
    window.appConfigLoaded = true; // Mark as loaded even on error
  }
})();

// Helper function for WASM to check if config is ready
window.isAppConfigReady = function() {
  return window.appConfigLoaded;
};

// Helper function for WASM to get the config
window.getAppConfig = function() {
  if (!window.appConfigLoaded) {
    console.warn("Config not yet loaded");
    return null;
  }
  if (window.appConfigError) {
    console.error("Config failed to load:", window.appConfigError);
    return null;
  }
  return window.appConfig;
};

console.log("common.js loaded - app config fetch initiated");
