module.exports = {
  root: true,
  env: {
    node: true,
    mocha: true
  },
  extends: ['standard'],
  overrides: [
    {
      files: ['*.test.js'],
      rules: {
        'no-unused-expressions': 'off'
      }
    }
  ]
}
