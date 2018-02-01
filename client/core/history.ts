// history types: github.com/DefinitelyTyped/DefinitelyTyped/issues/14062
import { createHistory as createHistoryFn } from 'history';
import { useQueries } from 'history';
const qs = require('qs')

export const createHistory = function() {
  return useQueries(createHistoryFn)({
    stringifyQuery(ob) {
      return qs.stringify(ob, {encode: true})
    },
    parseQueryString(str) {
      return qs.parse(str)
    }
  })
}
