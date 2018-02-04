import { createMemoryHistory, History } from 'history';
// https://github.com/ReactTraining/react-router/issues/2907
var createBrowserHistory = require('history/lib/createBrowserHistory');

const qs = require('qs')

export const encodeObjUrl = (obj) => {
  return qs.stringify(obj, { encode: true, encodeValuesOnly: true})
}

export const decodeObjString = (str) => {
  return qs.parse(str)
}

export const supportsHistory = () => {
  return typeof window !== 'undefined' && !!window.history
}

export const createHistoryInstance = function(): History {
  return supportsHistory() ? createBrowserHistory(): createMemoryHistory()
}
