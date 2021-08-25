'use strict'

const axios = require('axios')
const chardet = require('chardet')
const iconv = require('iconv-lite')

module.exports.execute = (event, context, callback) => {
  const url = decodeURIComponent(event.queryStringParameters.url)
  axios.get(url, {
    responseType: 'arraybuffer',
    transformResponse: data => {
      const encoding = chardet.detect(data)
      if (!encoding) {
        throw new Error('chardet failed to detect encoding')
      }
      return iconv.decode(data, encoding)
    }
  }).then((res) => {
    const result = /<title>(.+)<\/title>/.exec(res.data)
    const title = result ? result[1] : ''
    const response = {
      statusCode: 200,
      body: JSON.stringify({ title })
    }
    callback(null, response)
  })
}
