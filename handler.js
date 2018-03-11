'use strict';

const axios = require('axios');

module.exports.getWebTitle = (event, context, callback) => {
  const url = event.queryStringParameters.url;
  axios.get(url).then((e) => {
    const response = {
      statusCode: 200,
      body: JSON.stringify(e),
    };
    callback(null, response);
  });
};
