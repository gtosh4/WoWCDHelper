module.exports = {
  parser: "vue-eslint-parser",
  "parserOptions": {
      "parser": "@typescript-eslint/parser",
      "sourceType": "module",
  },
  extends: [
    "eslint:recommended",
    'plugin:vue/recommended',
  ],
  rules: {
    "semi": ["error", "never"],
    "indent": ["error", 2],
    "vue/max-attributes-per-line": ["error", {
      "singleline": 3,
      "multiline": {
        "max": 1,
        "allowFirstLine": false
      }
    }],
    "no-console": ["error", { allow: ["warn", "error"] }],
    "vue/multiline-html-element-content-newline": ["error", {
      allowEmptyLines: true,
    }],
    "no-unused-vars": ["error", {
      args: "none"
    }],
  }
}
