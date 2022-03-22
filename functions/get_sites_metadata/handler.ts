import { APIGatewayProxyHandler } from 'aws-lambda'
import 'source-map-support/register'
import axios from 'axios'
import * as chardet from 'chardet'
import * as iconv from 'iconv-lite'

export const execute: APIGatewayProxyHandler = (event, context, callback) => {
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
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Credentials': true,
      },
      body: JSON.stringify({ title })
    }
    callback(null, response)
  })
}
