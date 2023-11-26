import { middyfy } from '~/core/middyfy'
import { createHash } from 'node:crypto'
import * as AWS from 'aws-sdk'
const s3 = new AWS.S3()

const lambdaHandler = async (event) => {
  const body = event.body
  const hash = createHash('md5').update(body.attachment)
  // 先頭の ~;base64, まではファイルデータとして不要なので空文字で置換する
  const fileData = body.attachment.replace(/^data:\w+\/\w+;base64,/, '')
  const decodedFile = Buffer.from(fileData, 'base64')
  const bucket = 'melt-storage'
  const key = `images/${hash.digest('hex')}.${body.key.split('.')[1]}`
  const url = `https://s3-ap-northeast-1.amazonaws.com/${bucket}/${key}`
  const params = {
    Bucket: bucket,
    Key: key,
    ContentType: body.type,
    Body: decodedFile,
  }
  await s3.putObject(params).promise()
  return {
    statusCode: 200,
    body: JSON.stringify({
      name: body.key,
      url: url,
    }),
  }
}

export const main = middyfy(lambdaHandler)
