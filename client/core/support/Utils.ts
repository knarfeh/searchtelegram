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

  static collapse(collection, seed) {
    const reducer = (current, fn)=> fn(current)
    return reduce(collection, reducer, seed)
  }

  static instanceOf(klass) {
    return (val)=> val instanceof klass
  }

  static interpolate(str, interpolations){
    return str.replace(
  		/{([^{}]*)}/g,
  		(a, b) => {
  			var r = interpolations[b];
  			return typeof r === 'string' || typeof r === 'number' ? r : a;
  		}
    )
  }

  static translate(key, interpolations?) {
    if (interpolations) {
      return Utils.interpolate(key, interpolations)
    } else {
      return key
    }
  }

}
