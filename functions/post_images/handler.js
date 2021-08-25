'use strict'

const AWS = require('aws-sdk')
const s3 = new AWS.S3()
const crypto = require('crypto')

module.exports.execute = (event, context, callback) => {
  const body = JSON.parse(event.body)
  const hash = crypto.createHash('md5').update(body.attachment)
  // 先頭の ~;base64, まではファイルデータとして不要なので空文字で置換する
  const fileData = body.attachment.replace(/^data:\w+\/\w+;base64,/, '')
  const decodedFile = Buffer.from(fileData, 'base64')
  const bucket = 'melt-storage'
  const key = `images/${hash.digest('hex')}.${body.key.split('.')[1]}`
  const url = `https://s3-ap-northeast-1.amazonaws.com/${bucket}/${key}`
  s3.putObject({
    Bucket: bucket,
    Key: key,
    ContentType: body.type,
    Body: decodedFile
  }, function (err, data) {
    if (err) {
      callback(err)
    } else {
      const response = {
        statusCode: 200,
        body: JSON.stringify({
          name: body.key,
          url: url
        })
      }
      callback(null, response)
    }
  })
}
