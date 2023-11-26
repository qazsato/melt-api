import { middyfy } from '~/core/middyfy'
import axios from 'axios'

const lambdaHandler = async () => {
  const url = 'https://api.github.com/repos/qazsato/melt/releases/latest'
  const res = await axios.get(url)
  const body = {
    version: res.data.tag_name,
  }
  return {
    statusCode: 200,
    body: JSON.stringify(body),
  }
}

export const main = middyfy(lambdaHandler)
