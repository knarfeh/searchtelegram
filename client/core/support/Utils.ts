const reduce = require("lodash/reduce")
const map = require("lodash/map")
const reject = require("lodash/reject")
const isUndefined = require("lodash/isUndefined")

export class Utils {
  static guidCounter = 0

  static guid(prefix=""){
    let id = ++Utils.guidCounter
    return prefix.toString() + id
  }

  
}