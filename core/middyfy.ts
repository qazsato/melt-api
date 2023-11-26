import middy from '@middy/core'
import { apiAuthMiddleware } from '~/middlewares/apiAuth'
import httpErrorHandler from '@middy/http-error-handler'
import middyJsonBodyParser from '@middy/http-json-body-parser'

export const middyfy = (handler) => {
  return middy(handler)
    .use(apiAuthMiddleware)
    .use(httpErrorHandler())
    .use(middyJsonBodyParser())
}
