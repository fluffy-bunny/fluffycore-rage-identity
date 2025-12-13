/**
 * RAGE Logging Helper
 * 
 * This script provides pre-WASM logging control via localStorage.
 * Use this before the WASM module loads to enable logging from the start.
 */

(function() {
    'use strict';

    // Create the rage namespace early
    window.rage = window.rage || {};

    /**
     * Set the logging level (persists to localStorage)
     * This can be called before WASM loads
     * @param {string} level - Log level: 'disabled', 'trace', 'debug', 'info', 'warn', 'error', 'fatal', 'panic'
     */
    window.rage.SetLogLevel = function(level) {
        var validLevels = ['disabled', 'trace', 'debug', 'info', 'warn', 'error', 'fatal', 'panic'];
        var normalizedLevel = level.toLowerCase();
        
        if (validLevels.indexOf(normalizedLevel) === -1) {
            console.error('[RAGE] Invalid log level:', level, '- Valid levels:', validLevels.join(', '));
            return;
        }
        
        localStorage.setItem('rage_logging_level', normalizedLevel);
        console.log('[RAGE] Logging level saved:', normalizedLevel);
        console.log('[RAGE] Reload the page for changes to take effect (or wait for WASM to load)');
    };

    /**
     * Enable or disable logging (persists to localStorage)
     * This is a convenience method that sets level to 'debug' or 'disabled'
     * @param {boolean} enabled - true to enable logging (debug level), false to disable
     */
    window.rage.EnableLogging = function(enabled) {
        window.rage.SetLogLevel(enabled ? 'debug' : 'disabled');
    };

    /**
     * Get the current logging level
     * @returns {string} Current log level or 'disabled' if not set
     */
    window.rage.GetLogLevel = function() {
        return localStorage.getItem('rage_logging_level') || 'disabled';
    };

    /**
     * Check if logging is currently enabled (not disabled)
     * @returns {boolean} true if logging is enabled
     */
    window.rage.IsLoggingEnabled = function() {
        var level = window.rage.GetLogLevel();
        return level !== 'disabled';
    };

    /**
     * Clear the logging preference (resets to default disabled state)
     */
    window.rage.ClearLoggingPreference = function() {
        localStorage.removeItem('rage_logging_level');
        console.log('[RAGE] Logging preference cleared. Reload the page to reset to default (disabled).');
    };

    // Log availability
    var currentLevel = window.rage.GetLogLevel();
    console.log('[RAGE] Pre-WASM logging helper loaded. Current level:', currentLevel);
    console.log('[RAGE] Use rage.SetLogLevel("debug"|"info"|"warn"|"error"|"disabled") to change level');
    console.log('[RAGE] Or use rage.EnableLogging(true/false) for simple on/off');
})();
