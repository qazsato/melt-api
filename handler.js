'use strict';

const axios = require('axios');

module.exports.getWebTitle = (event, context, callback) => {
  const url = decodeURIComponent(event.queryStringParameters.url);
  axios.get(url).then((e) => {
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
  const response = {
    statusCode: 200,
    body: data.content
  };
  callback(null, response);
};
