'use strict';

const AWS = require('aws-sdk');
const s3 = new AWS.S3();
const axios = require('axios');
const marked = require('marked');
const emoji = require('node-emoji');
const highlight = require('highlight.js');
const chardet = require('chardet')
const iconv = require('iconv-lite')

module.exports.getWebTitle = (event, context, callback) => {
  const url = decodeURIComponent(event.queryStringParameters.url);
  axios.get(url, {
    responseType: 'arraybuffer',
    transformResponse: data => {
      const encoding = chardet.detect(data);
      if (!encoding) {
        throw new Error('chardet failed to detect encoding');
      }
      return iconv.decode(data, encoding);
    }
  }).then((e) => {
    const title = /<title>(.+)<\/title>/.exec(e.data)[1];
    const response = {
      statusCode: 200,
      body: JSON.stringify({title}),
    };
    callback(null, response);
  });
};

module.exports.postNote = (event, context, callback) => {
  const data = JSON.parse(event.body);
  const content = data.content;
  const file = /.+\/(.+)\.json/.exec(data.path)[1];
  s3.putObject({
    Bucket: 'melt-storage',
    Key: `note/${file}.html`,
    Body: generateNote(content),
    ContentType: 'text/html'
  }, function (err, data) {
    if (err) {
      callback(err);
    } else {
      const response = {
        statusCode: 200,
        body: JSON.stringify({
          url: `https://s3-ap-northeast-1.amazonaws.com/melt-storage/note/${file}.html`
        })
      };
      callback(null, response);
    }
  });
};

function generateNote(content) {
  marked.setOptions({
    renderer: createRenderer(),
    highlight: (code, lang) => {
      if (lang) {
        return highlight.highlightAuto(code, [lang]).value;
      }
      return code;
    }
  });
  return `
    <!DOCTYPE html>
    <html lang="ja">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <meta http-equiv="X-UA-Compatible" content="ie=edge">
      <title>Document</title>
      <link rel="stylesheet" href="./lib/github-markdown.css">
      <link rel="stylesheet" href="./lib/github-highlight.css">
      <style>
        .markdown-body {
          padding: 15px 30px;
        }
        .markdown-body .check-list {
          list-style: none;
        }
        .markdown-body .check-list input[type="checkbox"] {
          margin-left: -1.3em;
          margin-right: 0.2em;
        }
      </style>
    </head>
    <body>
      <article class="markdown-body">
        ${marked(content)}
      </article>
    </body>
    </html>
  `;
}

function createRenderer() {
  const renderer = new marked.Renderer();
  renderer.list = (body, ordered) => {
    // GFMのCheckbox記法に対応するため拡張
    let html;
    if (ordered === true) {
      html = `<ol>${body}</ol>`;
    } else if (body.includes('type="checkbox"')) {
      html = `<ul class="check-list">${body}</ul>`
    } else {
      html = `<ul>${body}</ul>`;
    }
    return html;
  };
  renderer.listitem = (text) => {
    // GFMのCheckbox記法に対応するため拡張
    text = text.replace(/\[x\]/gi, '<input type="checkbox" disabled checked>');
    text = text.replace(/\[ \]/gi, '<input type="checkbox" disabled>');
    return `<li>${text}</li>`;
  };
  renderer.text = (html) => {
    // 絵文字に対応するため拡張
    return emoji.emojify(html);
  };
  renderer.link = (href, title, text) => {
    // 新規タブでブラウザ起動するため拡張
    return `<a target="_blank" href="${href}" title="${title}">${text}</a>`
  };
  return renderer;
}
