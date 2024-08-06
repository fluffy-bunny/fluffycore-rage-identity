const simpleImportSort = require('eslint-plugin-simple-import-sort');
const prettierPlugin = require('eslint-plugin-prettier');
const prettierConfig = require('eslint-config-prettier');
const typescriptEslintParser = require('@typescript-eslint/parser');
const typescriptEslintPlugin = require('@typescript-eslint/eslint-plugin');

module.exports = [
  {
    files: ['src/**/*.{js,jsx,ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2021,
      sourceType: 'module',
      parser: typescriptEslintParser,
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
        tsconfigRootDir: __dirname,
        project: ['./tsconfig.json'], // Adjust this path if your `tsconfig.json` is elsewhere
      },
      globals: {
        browser: true,
        es2021: true,
      },
    },
    plugins: {
      'simple-import-sort': simpleImportSort,
      prettier: prettierPlugin,
      '@typescript-eslint': typescriptEslintPlugin,
    },
    rules: {
      'simple-import-sort/imports': 'error',
      'simple-import-sort/exports': 'error',
      'prettier/prettier': 'error',
      '@typescript-eslint/explicit-module-boundary-types': 'off', // Adjust TypeScript rules as needed
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_' },
      ],
    },
  },
  {
    // Apply prettier config
    rules: prettierConfig.rules,
  },
];
