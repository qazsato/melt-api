import createError from 'http-errors'
import AWS from 'aws-sdk'
import { GetItemInput } from 'aws-sdk/clients/dynamodb'

export const apiAuthMiddleware = {
  before: async (request) => {
    const apiKey = request.event.queryStringParameters?.api_key
    if (!apiKey) {
      const message = 'api_key の指定が必要です'
      throw new createError.Unauthorized(JSON.stringify({ message }))
    }
    const isExist = await isExistKey(apiKey)
    if (!isExist) {
      const message = 'api_key が不正です'
      throw new createError.Unauthorized(JSON.stringify({ message }))
    }
  },
}

const isExistKey = async (key: string): Promise<boolean> => {
  const dynamodb = new AWS.DynamoDB()
  const tableName = process.env.KEY_TABLE_NAME
  if (!tableName) {
    throw new Error('KEY_TABLE_NAME environment variable is not set')
  }
  const params: GetItemInput = {
    TableName: tableName,
    Key: {
      key: { S: key },
    },
  }
  const result = await dynamodb.getItem(params).promise()
  if (!result.Item || !result.Item['key']) {
    return false
  }
  return true
}
