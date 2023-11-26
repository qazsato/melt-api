import { middyfy } from '~/core/middyfy'
import axios from 'axios'
import chardet from 'chardet'
import iconv from 'iconv-lite'

const lambdaHandler = async (event) => {
  const url = decodeURIComponent(event.queryStringParameters.url)
  const res = await axios.get(url, {
    responseType: 'arraybuffer',
    transformResponse: (data) => {
      const encoding = chardet.detect(data)
      if (!encoding) {
        throw new Error('chardet failed to detect encoding')
      }
      return iconv.decode(data, encoding)
    },
  })
  const result = /<title>(.+)<\/title>/.exec(res.data)
  const title = result ? result[1] : ''
  return {
    statusCode: 200,
    body: JSON.stringify({ title }),
  }
}

export const main = middyfy(lambdaHandler)
